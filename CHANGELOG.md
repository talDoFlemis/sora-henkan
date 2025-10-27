## [1.21.0](https://github.com/talDoFlemis/sora-henkan/compare/v1.20.0...v1.21.0) (2025-10-27)

### Features

* add swagger ([8fe41d3](https://github.com/talDoFlemis/sora-henkan/commit/8fe41d3a8ae8f32e71400beec7c5987160ef0667))

## [1.20.0](https://github.com/talDoFlemis/sora-henkan/compare/v1.19.0...v1.20.0) (2025-10-27)

### Features

* **ec2-otel:** wait for nat gateway before creating ([f123b4a](https://github.com/talDoFlemis/sora-henkan/commit/f123b4a3ba0e7665d4affc7ce74bc2e65048e192))
* **user_data_otel:** always rerun script on startup ([ed3403c](https://github.com/talDoFlemis/sora-henkan/commit/ed3403c74216d7c1de538dc77a99b1164136bbf2))

### Bug Fixes

* **ec2_otel:** remove awsemf to reduce costs ([4df37e2](https://github.com/talDoFlemis/sora-henkan/commit/4df37e2c68b162943cf9b88a77ff323ceb64c6c2))

## [1.19.0](https://github.com/talDoFlemis/sora-henkan/compare/v1.18.0...v1.19.0) (2025-10-24)

### Features

* **ami_builder_data:** add always run cloud init and containers limits and requests ([7f93d21](https://github.com/talDoFlemis/sora-henkan/commit/7f93d21e98308d1d054bc175c133e31ae9af7241))

## [1.18.0](https://github.com/talDoFlemis/sora-henkan/compare/v1.17.0...v1.18.0) (2025-10-24)

### Features

* expose worker health endpoint on alb ([0647ab1](https://github.com/talDoFlemis/sora-henkan/commit/0647ab16b44f1c9b0ed2b9bcf4d14d0b3fd90fb1))

## [1.17.0](https://github.com/talDoFlemis/sora-henkan/compare/v1.16.2...v1.17.0) (2025-10-24)

### Features

* add worker healthcheck endpoint to alb, sg e user_data to handle scaling if down ([646dfbd](https://github.com/talDoFlemis/sora-henkan/commit/646dfbdad431e3cac9694ec3a008aab9adc22fe0))

## [1.16.2](https://github.com/talDoFlemis/sora-henkan/compare/v1.16.1...v1.16.2) (2025-10-24)

### Bug Fixes

* **script.js:** use inline urls ([56bc8e7](https://github.com/talDoFlemis/sora-henkan/commit/56bc8e71fa617b06cb793531357ef0632752b9e5))

## [1.16.1](https://github.com/talDoFlemis/sora-henkan/compare/v1.16.0...v1.16.1) (2025-10-24)

### Bug Fixes

* **user_data_ami_builder:** just do a docker compose pull ([b5cbd63](https://github.com/talDoFlemis/sora-henkan/commit/b5cbd6342f274ee42920ecb54d3ab8d8dba3a531))

## [1.16.0](https://github.com/talDoFlemis/sora-henkan/compare/v1.15.0...v1.16.0) (2025-10-24)

### Features

* add ami image builder ([01f9832](https://github.com/talDoFlemis/sora-henkan/commit/01f9832b3e6d878ae3759c1ae49635a9e01f6277))
* use ami to launch new instances ([ad22f7e](https://github.com/talDoFlemis/sora-henkan/commit/ad22f7e09f68d8f3d3af6b7ddd91f8625a546d52))

## [1.15.0](https://github.com/talDoFlemis/sora-henkan/compare/v1.14.0...v1.15.0) (2025-10-24)

### Features

* add animes.json images ([9a0e8b2](https://github.com/talDoFlemis/sora-henkan/commit/9a0e8b23dd49741555eb77ee265f129acfd1ebd9))
* add load test script ([08288fe](https://github.com/talDoFlemis/sora-henkan/commit/08288fee75b8a8e750007d13077177d2373a80ef))

### Bug Fixes

* **user_data_otel:** change logging to debug ([76fe50e](https://github.com/talDoFlemis/sora-henkan/commit/76fe50e628be8084564496b9f3a1231ebe95a4f7))
* **user_data:** missing -y to cloudwatch agent ([80afaed](https://github.com/talDoFlemis/sora-henkan/commit/80afaed0169b615bcc7c7b2c4493fdf298e75cb3))

## [1.14.0](https://github.com/talDoFlemis/sora-henkan/compare/v1.13.0...v1.14.0) (2025-10-24)

### Features

* add cloudflare client ip extractor ([92d98da](https://github.com/talDoFlemis/sora-henkan/commit/92d98daf082535c142c36e0b9cea8d515e4b1f60))

## [1.13.0](https://github.com/talDoFlemis/sora-henkan/compare/v1.12.0...v1.13.0) (2025-10-24)

### Features

* add cloudwatch-agent to ec2 instances ([7ea74b9](https://github.com/talDoFlemis/sora-henkan/commit/7ea74b91f002588780f23d446dea021e71d7b5a1))
* add target group attachment of jaeger to alb ([0e50a37](https://github.com/talDoFlemis/sora-henkan/commit/0e50a3791e100a8a06c9f3f1691d1adc9fb27a6d))
* **sg:** add ingress for jaeger ui from LB ([37cc333](https://github.com/talDoFlemis/sora-henkan/commit/37cc3339677b7dbf02c47f06533c935e30096b18))

### Bug Fixes

* **autoscaling:** remove jaeger from autoscaling group ([d9402f3](https://github.com/talDoFlemis/sora-henkan/commit/d9402f3f2d3fa4651f073e3a9657ee946c2092b8))

## [1.12.0](https://github.com/talDoFlemis/sora-henkan/compare/v1.11.1...v1.12.0) (2025-10-24)

### Features

* **sg:** add jaeger ui to alb security group ([b9b1f5b](https://github.com/talDoFlemis/sora-henkan/commit/b9b1f5b4cc6958de767eb466dacca34f7bc183d4))

## [1.11.1](https://github.com/talDoFlemis/sora-henkan/compare/v1.11.0...v1.11.1) (2025-10-24)

### Bug Fixes

* **user_data_otel:** bad jaeger image ([7a6b592](https://github.com/talDoFlemis/sora-henkan/commit/7a6b5929fea20209f17b97f0cfc6a60ee0b114f1))

## [1.11.0](https://github.com/talDoFlemis/sora-henkan/compare/v1.10.0...v1.11.0) (2025-10-24)

### Features

* add jaeger lb url to alb ([cd73a54](https://github.com/talDoFlemis/sora-henkan/commit/cd73a545d3e933bd7d383132e6fb1160c3c01c62))
* **s3:** make bucket public ([79bd830](https://github.com/talDoFlemis/sora-henkan/commit/79bd830896e9347d60ea2a9de746f20158e4219a))

## [1.10.0](https://github.com/talDoFlemis/sora-henkan/compare/v1.9.0...v1.10.0) (2025-10-24)

### Features

* add jaeger target group ([01521ca](https://github.com/talDoFlemis/sora-henkan/commit/01521cacf0b6e9f92c2ec3cd3c11b3da7e0ce915))
* **router:** add logger after running otel thing ([74385fd](https://github.com/talDoFlemis/sora-henkan/commit/74385fd356936781ee7d164018b1fe87b8113b64))
* **s3:** add cors policy to allow all origins ([9d40820](https://github.com/talDoFlemis/sora-henkan/commit/9d40820a0a94f8857acd55acaffe846ead4a0bee))
* **user_data_otel:** add jaeger ([9bc1d53](https://github.com/talDoFlemis/sora-henkan/commit/9bc1d53530418bdec3b1c637b535e499fd7e52cc))

## [1.9.0](https://github.com/talDoFlemis/sora-henkan/compare/v1.8.0...v1.9.0) (2025-10-24)

### Features

* add mime type extension to save image ([60bb1ad](https://github.com/talDoFlemis/sora-henkan/commit/60bb1ade43c8403d65142ef48e0e830c39ca82e4))

## [1.8.0](https://github.com/talDoFlemis/sora-henkan/compare/v1.7.1...v1.8.0) (2025-10-24)

### Features

* **image_usecase:** add bucket name log ([5fa1db6](https://github.com/talDoFlemis/sora-henkan/commit/5fa1db6269bf1b6559452e8d367afab38a37bc50))

## [1.7.1](https://github.com/talDoFlemis/sora-henkan/compare/v1.7.0...v1.7.1) (2025-10-24)

### Bug Fixes

* add image processor bucket name variable to interpolate ([9d4afc3](https://github.com/talDoFlemis/sora-henkan/commit/9d4afc36eb8254b6c5064d22e518803f361e66eb))

## [1.7.0](https://github.com/talDoFlemis/sora-henkan/compare/v1.6.1...v1.7.0) (2025-10-24)

### Features

* add CI/CD to frontend and go apps ([d9c9540](https://github.com/talDoFlemis/sora-henkan/commit/d9c9540b048d3f824e2d50a2940b79c11cdbdb21))
* add format to frontend ([521d9eb](https://github.com/talDoFlemis/sora-henkan/commit/521d9eb5ab71209051c033d7f448f25880490231))

### Documentation

* add README ([807e4d5](https://github.com/talDoFlemis/sora-henkan/commit/807e4d5fbb5733d5f3431254325533b17417988e))

## [1.6.1](https://github.com/talDoFlemis/sora-henkan/compare/v1.6.0...v1.6.1) (2025-10-24)

### Bug Fixes

* duple interpolation ([7ee5a3a](https://github.com/talDoFlemis/sora-henkan/commit/7ee5a3ab617e3a33b04b3a2d652e9e52ba5d48be))

## [1.6.0](https://github.com/talDoFlemis/sora-henkan/compare/v1.5.0...v1.6.0) (2025-10-24)

### Features

* inject vite api url ([eabdbc7](https://github.com/talDoFlemis/sora-henkan/commit/eabdbc7716f66b3351b3f5667132449b70a4122d))

## [1.5.0](https://github.com/talDoFlemis/sora-henkan/compare/v1.4.7...v1.5.0) (2025-10-24)

### Features

* change namespace and app names ([a195cb7](https://github.com/talDoFlemis/sora-henkan/commit/a195cb7869f98437385a11b89dae749ee8658a80))

### Bug Fixes

* add aws region resolution ([3da2376](https://github.com/talDoFlemis/sora-henkan/commit/3da237610bb569acdc0e2a867b14c64a955e7797))

## [1.4.7](https://github.com/talDoFlemis/sora-henkan/compare/v1.4.6...v1.4.7) (2025-10-24)

### Bug Fixes

* pad access and secret template ([1fc08db](https://github.com/talDoFlemis/sora-henkan/commit/1fc08db885143349a5c77b6411db1ea9123227c2))

## [1.4.6](https://github.com/talDoFlemis/sora-henkan/compare/v1.4.5...v1.4.6) (2025-10-24)

### Bug Fixes

* **settings:** disable anonymous ([25a32b4](https://github.com/talDoFlemis/sora-henkan/commit/25a32b435c5b19fd4393e1f2f0792f5894354e69))

## [1.4.5](https://github.com/talDoFlemis/sora-henkan/compare/v1.4.4...v1.4.5) (2025-10-24)

### Bug Fixes

* **alb:** health endpoint for api is wrong ([f079be9](https://github.com/talDoFlemis/sora-henkan/commit/f079be95cbcf74945422ad9ca748f2b6885424fa))

## [1.4.4](https://github.com/talDoFlemis/sora-henkan/compare/v1.4.3...v1.4.4) (2025-10-24)

### Bug Fixes

* enabled opentelemetry on ec2 ([a00329d](https://github.com/talDoFlemis/sora-henkan/commit/a00329dd0f09020d6f8eaf96f48a477caafd1416))

## [1.4.3](https://github.com/talDoFlemis/sora-henkan/compare/v1.4.2...v1.4.3) (2025-10-24)

### Bug Fixes

* **ec2:** make host mode for handling IAM minio access ([59c66ae](https://github.com/talDoFlemis/sora-henkan/commit/59c66aee811b8213b08ee935f88a0e71e8c9a694))

## [1.4.2](https://github.com/talDoFlemis/sora-henkan/compare/v1.4.1...v1.4.2) (2025-10-24)

### Bug Fixes

* **settings:** use NewIAM if on aws s3 in minio client ([90b967f](https://github.com/talDoFlemis/sora-henkan/commit/90b967f6bbe8758cd55a8c8cf2244a5999287e17))

## [1.4.1](https://github.com/talDoFlemis/sora-henkan/compare/v1.4.0...v1.4.1) (2025-10-24)

### Bug Fixes

* **settings:** cannot interpolat e object storage env variables ([a869413](https://github.com/talDoFlemis/sora-henkan/commit/a869413f2605b69d26bab98a01a8b829fdbd7e92))

## [1.4.0](https://github.com/talDoFlemis/sora-henkan/compare/v1.3.0...v1.4.0) (2025-10-23)

### Features

* add frontend ([d3b145b](https://github.com/talDoFlemis/sora-henkan/commit/d3b145b65e405f4e1eee469d5fb33c7d918ee01b))
* add frontend to ec2 ([e74a1ac](https://github.com/talDoFlemis/sora-henkan/commit/e74a1ac2cf05ff98d5b1ba74c304a1d482d20ee9))

## [1.3.0](https://github.com/talDoFlemis/sora-henkan/compare/v1.2.0...v1.3.0) (2025-10-23)

### Features

* add handler per port and port ([83ee2b3](https://github.com/talDoFlemis/sora-henkan/commit/83ee2b3e33d31ac6d9ebf0cb977157adb0cd56ba))

### Bug Fixes

* **settings:** ssl mode using underscore ([ad56c87](https://github.com/talDoFlemis/sora-henkan/commit/ad56c87168e99747e6d04e5d06ba999d8cb27e3f))

## [1.2.0](https://github.com/talDoFlemis/sora-henkan/compare/v1.1.0...v1.2.0) (2025-10-23)

### Features

* add base terraform ([82fd1ed](https://github.com/talDoFlemis/sora-henkan/commit/82fd1ed933fd2bc7a2abf99b11b3a3c21fe6f2dd))
* add lb between frontenzo and back ([df68876](https://github.com/talDoFlemis/sora-henkan/commit/df68876d61cad4028a840d8c5b91ecdd922c0c06))
* add migrate dockerfile ([6a99b13](https://github.com/talDoFlemis/sora-henkan/commit/6a99b138afaaa100afdc9fd4ce6cf05c00a6c406))

### Bug Fixes

* add require ssl mode for rds deploy ([9e940c6](https://github.com/talDoFlemis/sora-henkan/commit/9e940c6ebd56e618570c53ec5310eea1078bab3d))
* bad templanting and otel collector with bad logs to cloudwatch ([eb5dcb3](https://github.com/talDoFlemis/sora-henkan/commit/eb5dcb3c839f5cf0fc17f201e81719a5801ea1a7))

## [1.1.0](https://github.com/talDoFlemis/sora-henkan/compare/v1.0.0...v1.1.0) (2025-10-22)

### Features

* add docker image publishing to releases ([da697c4](https://github.com/talDoFlemis/sora-henkan/commit/da697c450b269eeb7685c6eb71440b6f35c6b366))

## 1.0.0 (2025-10-22)

### Features

* add air ([d231b5c](https://github.com/talDoFlemis/sora-henkan/commit/d231b5cccc5ac2694e862c9be8101798db2f3d00))
* add api with healthcheck handler ([6a48430](https://github.com/talDoFlemis/sora-henkan/commit/6a48430a6802b88056af7e6d6f976e69daf2ab6c))
* add Dockerfile ([c307816](https://github.com/talDoFlemis/sora-henkan/commit/c307816f513c322cdb661fa9f96ee90135b5d298))
* add image handler ([f5f3d30](https://github.com/talDoFlemis/sora-henkan/commit/f5f3d30e6383f158bbaa36689d59ad87e14e7f29))
* add integration between publisher and sub ([3a6ca52](https://github.com/talDoFlemis/sora-henkan/commit/3a6ca520fbeb8b25327f61faf05705e6dc717ef7))
* add localstack ([9ea0bca](https://github.com/talDoFlemis/sora-henkan/commit/9ea0bca582836ca4e80a7b0b022d830d4c8334f7))
* add migrate script ([4b96c6f](https://github.com/talDoFlemis/sora-henkan/commit/4b96c6f565101b359ff1af445a70e29c99a8cf82))
* add minio object storer ([45580d4](https://github.com/talDoFlemis/sora-henkan/commit/45580d4c03b5e2cf87636772f44df04123ee79a4))
* add scaling to lanczos too ([2a80f50](https://github.com/talDoFlemis/sora-henkan/commit/2a80f508bd0a39bad7dc37aa9b0b8d45ad786797))
* add semantic releaser ([7d2373e](https://github.com/talDoFlemis/sora-henkan/commit/7d2373e7c11e37734f3069b0b39734d8a0155ebf))
* add telemetry package ([eff28b5](https://github.com/talDoFlemis/sora-henkan/commit/eff28b590e48fa7d3de969120f03db9ba360b79d))
* add validation ([516e781](https://github.com/talDoFlemis/sora-henkan/commit/516e781017e3bf38f59a698b9d62e8818ccb77ce))
* add vips image processor ([3d05d2c](https://github.com/talDoFlemis/sora-henkan/commit/3d05d2c619b58a65cade6aee1019ef1aa114b13a))
* add worker ([43669b8](https://github.com/talDoFlemis/sora-henkan/commit/43669b8b79cb2f569696d438a7c9771c1c9b1024))
