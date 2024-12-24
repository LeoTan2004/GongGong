from unittest import TestCase

from xtu_ems.util.captcha import ImageDetector


class TestImageDetector(TestCase):
    def test_verify(self):
        """测试验证码识别"""
        with open('verifycode.jpeg', 'rb') as f:
            img = f.read()
        res = ImageDetector().verify(img)
        self.assertEqual('cnbz', res)
