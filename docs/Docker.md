# Minio Docker Quickstart Guide [![Gitter](https://badges.gitter.im/Join%20Chat.svg)](https://gitter.im/minio/minio?utm_source=badge&utm_medium=badge&utm_campaign=pr-badge&utm_content=badge)


## 1. Test Minio Docker Container
Minio generates new access and secret keys each time you run this command. Container state is lost after you end this session. This mode is only intended for testing purpose.

```bash
docker run -p 9000:9000 minio/minio /export
```

## 2. Run Minio Docker Container
Minio container requires a persistent volume to store configuration and application data. Following command maps local persistent directories from the host OS to virtual config `~/.minio` and export `/export` directories. 

```bash
docker run -p 9000:9000 --name minio1 \
  -v /mnt/export/minio1:/export \
  -v /mnt/config/minio1:/root/.minio \
  minio/minio /export
```

## 3. Custom Access and Secret Keys
To override Minio's auto-generated keys, you may pass secret and access keys explicitly as environment variables. Minio server also allows regular strings as access and secret keys.
```bash
docker run -p 9000:9000 --name minio1 \
  -e "MINIO_ACCESS_KEY=AKIAIOSFODNN7EXAMPLE" \
  -e "MINIO_SECRET_KEY=wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY" \
  -v /mnt/export/minio1:/export \
  -v /mnt/config/minio1:/root/.minio \
  minio/minio /export
```
