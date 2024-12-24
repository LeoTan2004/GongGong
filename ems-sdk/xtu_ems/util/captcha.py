"""验证码模块辅助通过工具"""
from abc import abstractmethod, ABC

import ddddocr


class Captcha(ABC):
    """验证码辅助校验器"""

    @abstractmethod
    def verify(self, target, *args, **kwargs) -> str:
        """验证码校验"""
        pass


_ocr = ddddocr.DdddOcr(show_ad=False)


class ImageDetector(Captcha):
    """图形识别验证码校验器"""

    def verify(self, target, *args, **kwargs) -> str:
        return _ocr.classification(img=target)
