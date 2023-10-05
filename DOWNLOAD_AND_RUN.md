# Download and run

## Direct download
[pgsectest](https://github.com/MannemSolutions/pgsectest) is available for download for many platforms and architectures from the [Github Releases page](https://github.com/MannemSolutions/pgsectest/releases).
It could be as simple as:
```bash
PGTESTER_VERSION=v0.3.0
cd $(mktemp -d)
curl -Lo "pgsectest-${PGTESTER_VERSION}-linux-amd64.tar.gz" "https://github.com/MannemSolutions/pgsectest/releases/download/${PGTESTER_VERSION}/pgsectest-${PGTESTER_VERSION}-linux-amd64.tar.gz"
tar -xvf "./pgsectest-${PGTESTER_VERSION}-linux-amd64.tar.gz"
mv pgsectest /usr/local/bin
cd -
```
After that you can run pgsectest directly from the prompt:
```bash
pgsectest ./mytest1.yml mytest2.yml
```
Or using stdin:
```bash
cat ./mytests*.yml | pgsectest
```

## Container image
For container environments [pgsectest](https://github.com/MannemSolutions/pgsectest) is also available on [dockerhub](https://hub.docker.com/repository/docker/mannemsolutions/pgsectest).
You can easily pull it with:
```bash
docker pull mannemsolutions/pgsectest
```

Using it would be as easy as:
```bash
cat testdata/pgsectest/tests.yaml | docker run -i mannemsolutions/pgsectest pgsectest
```

## docker-compose
You can use pgsectest with docker compose.
The docker-compose.yml file could have contents like this:
```yaml
services:
  pgsectest:
    image: mannemsolutions/pgsectest
    command: pgsectest /etc/pgtestdata/examples/tests1.yaml
  postgres:
    image: postgres:13
    environment:
      POSTGRES_HOST_AUTH_METHOD: 'md5'
      POSTGRES_PASSWORD: pgsectest
```

it could be as easy as:
```bash
docker-compose up
```

or with only output for pgsectest:
```bash
docker-compose up -d postgres
docker-compose up pgsectest
```

or with tests defined locally:
```bash
docker-compose up -d postgres
cat ./mytests*.yml | docker-compose up pgsectest
```

## Direct build

Although not advised, you can also directly build from source:
```bash
go install github.com/mannemsolutions/pgsectest/cmd/pgsectest@v0.3.0
```

After that you can run pgsectest directly from the prompt:
```bash
pgsectest ./mytest1.yml mytest2.yml
```

Or using stdin:
```bash
cat ./mytests*.yml | pgsectest
```
