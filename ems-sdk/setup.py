from setuptools import setup, find_packages

with open('../requirements.ems.txt', 'r') as f:
    requirements = [req.strip() for req in f.readlines()]
print(requirements)
setup(
    name='xtu_ems',
    version='0.0.1',
    author='LeoTan',
    author_email='leotan-studio@foxmail.com',
    packages=find_packages(''),
    classifiers=[
        "Programming Language :: Python :: 3",
        "License :: OSI Approved :: Apache License",
        "Operating System :: OS Independent",
    ],
    package_dir={'xtu_ems': 'xtu_ems'},
    python_requires='>=3.10',
    install_requires=requirements
)
