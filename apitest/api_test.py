#!/usr/bin/env python3
# -*- coding: utf-8 -*-

import requests
import json
import sys
from typing import Dict, Any, Union, Set
import logging
from datetime import datetime
import os

# 创建日志目录
log_dir = "logs"
if not os.path.exists(log_dir):
    os.makedirs(log_dir)

# 配置日志
logging.basicConfig(
    level=logging.INFO,
    format="%(asctime)s - %(levelname)s - %(message)s",
    handlers=[
        logging.FileHandler(
            os.path.join(
                log_dir, f'test_results_{datetime.now().strftime("%Y%m%d")}.log'
            )
        ),
        logging.StreamHandler(),
    ],
)
logger = logging.getLogger(__name__)


class APITester:
    # HTTP方法对应的成功状态码
    SUCCESS_STATUS_CODES = {
        "GET": {200},
        "POST": {200, 201},
        "PUT": {200, 201, 204},
        "DELETE": {200, 204},
        "PATCH": {200, 204},
    }

    def __init__(self, base_url: str = "http://localhost:38080"):
        self.base_url = base_url
        self.token = None
        self.headers = {"Content-Type": "application/json"}
        self.test_product_id = None
        self.test_order_id = None

    def update_auth_header(self):
        if self.token:
            self.headers["Authorization"] = f"Bearer {self.token}"

    def check_response(
        self,
        response: requests.Response,
        http_method: str = "GET",
        additional_status_codes: Set[int] = None,
    ) -> Dict[str, Any]:
        """
        检查API响应
        :param response: requests的响应对象
        :param http_method: HTTP方法（GET, POST, PUT, DELETE等）
        :param additional_status_codes: 额外允许的状态码集合
        :return: 响应的JSON数据
        """
        try:
            data = response.json()
            # 获取当前HTTP方法允许的成功状态码
            allowed_status_codes = self.SUCCESS_STATUS_CODES.get(
                http_method.upper(), {200}
            )

            # 添加额外允许的状态码
            if additional_status_codes:
                allowed_status_codes = allowed_status_codes.union(
                    additional_status_codes
                )

            if response.status_code not in allowed_status_codes:
                raise Exception(
                    f"Unexpected status code for {http_method}: {response.status_code}. "
                    f"Expected one of {allowed_status_codes}. Response: {data}"
                )

            # 检查业务状态码（如果存在）
            if isinstance(data, dict):
                code = data.get("code")
                if (
                    code is not None
                    and code != 200
                    and response.status_code not in (additional_status_codes or set())
                ):
                    raise Exception(f"Business logic error. Response: {data}")

            return data
        except json.JSONDecodeError:
            raise Exception(f"Invalid JSON response: {response.text}")

    def test_user_register(self) -> bool:
        logger.info("测试用户注册")
        try:
            response = requests.post(
                f"{self.base_url}/api/v1/users/register",
                headers=self.headers,
                json={
                    "username": "testuser",
                    "password": "123456",
                    "phone": "13800138000",
                },
            )
            # 对于用户注册，接受200、201和409（用户已存在）状态码
            self.check_response(response, "POST", additional_status_codes={409})
            logger.info("用户注册测试通过")
            return True
        except Exception as e:
            logger.error(f"用户注册测试失败: {str(e)}")
            return False

    def test_user_login(self) -> bool:
        logger.info("测试用户登录")
        try:
            response = requests.post(
                f"{self.base_url}/api/v1/users/login",
                headers=self.headers,
                json={"username": "testuser", "password": "123456"},
            )
            data = self.check_response(response, "POST")
            self.token = data["data"]["token"]
            self.update_auth_header()
            logger.info("用户登录测试通过")
            return True
        except Exception as e:
            logger.error(f"用户登录测试失败: {str(e)}")
            return False

    def test_user_profile(self) -> bool:
        logger.info("测试获取用户信息")
        try:
            response = requests.get(
                f"{self.base_url}/api/v1/users/profile", headers=self.headers
            )
            self.check_response(response, "GET")
            logger.info("获取用户信息测试通过")
            return True
        except Exception as e:
            logger.error(f"获取用户信息测试失败: {str(e)}")
            return False

    def test_user_profile_update(self) -> bool:
        logger.info("测试更新用户信息")
        try:
            response = requests.put(
                f"{self.base_url}/api/v1/users/profile",
                headers=self.headers,
                json={"phone": "13800138001", "address": "test address"},
            )
            self.check_response(response, "PUT")
            logger.info("更新用户信息测试通过")
            return True
        except Exception as e:
            logger.error(f"更新用户信息测试失败: {str(e)}")
            return False

    def test_product_operations(self) -> bool:
        logger.info("测试商品相关操作")
        try:
            # 创建商品
            response = requests.post(
                f"{self.base_url}/api/v1/products",
                headers=self.headers,
                json={
                    "name": "test product",
                    "description": "This is a test product",
                    "price": 9.99,
                    "stock": 100,
                },
            )
            data = self.check_response(response, "POST")
            self.test_product_id = data["data"]["id"]

            # 获取商品列表
            response = requests.get(
                f"{self.base_url}/api/v1/products?page=1&page_size=10",
                headers=self.headers,
            )
            self.check_response(response, "GET")

            # 获取商品详情
            response = requests.get(
                f"{self.base_url}/api/v1/products/{self.test_product_id}",
                headers=self.headers,
            )
            self.check_response(response, "GET")

            logger.info("商品相关操作测试通过")
            return True
        except Exception as e:
            logger.error(f"商品相关操作测试失败: {str(e)}")
            return False

    def test_order_operations(self) -> bool:
        logger.info("测试订单相关操作")
        try:
            # 创建订单
            response = requests.post(
                f"{self.base_url}/api/v1/orders",
                headers=self.headers,
                json={"product_id": self.test_product_id, "quantity": 1},
            )
            data = self.check_response(response, "POST")
            self.test_order_id = data["data"]["id"]

            # 获取订单列表
            response = requests.get(
                f"{self.base_url}/api/v1/orders?page=1&page_size=10",
                headers=self.headers,
            )
            self.check_response(response, "GET")

            # 获取订单详情
            response = requests.get(
                f"{self.base_url}/api/v1/orders/{self.test_order_id}",
                headers=self.headers,
            )
            self.check_response(response, "GET")

            # 获取订单区块链交易信息
            response = requests.get(
                f"{self.base_url}/api/v1/orders/{self.test_order_id}/transaction",
                headers=self.headers,
            )
            self.check_response(response, "GET")

            logger.info("订单相关操作测试通过")
            return True
        except Exception as e:
            logger.error(f"订单相关操作测试失败: {str(e)}")
            return False


def main():
    tester = APITester()

    # 定义测试用例列表
    test_cases = [
        (tester.test_user_register, "用户注册"),
        (tester.test_user_login, "用户登录"),
        (tester.test_user_profile, "获取用户信息"),
        (tester.test_user_profile_update, "更新用户信息"),
        (tester.test_product_operations, "商品相关操作"),
        (tester.test_order_operations, "订单相关操作"),
    ]

    # 执行所有测试用例
    all_passed = True
    for test_func, test_name in test_cases:
        logger.info(f"\n{'='*20} 开始测试: {test_name} {'='*20}")
        if not test_func():
            all_passed = False
            logger.error(f"{test_name} 测试失败")
            break
        logger.info(f"{'='*20} {test_name} 测试完成 {'='*20}\n")

    if all_passed:
        logger.info("所有API测试通过！")
        sys.exit(0)
    else:
        logger.error("API测试失败！")
        sys.exit(1)


if __name__ == "__main__":
    main()
