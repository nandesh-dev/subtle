# Subtle

![Cover](https://github.com/nandesh-dev/subtle/blob/main/doc/assets/images/cover.png?raw=true)

Your tool to manage all the subtitles in your personal media library. Subtle helps you extract and organize subtitle files for all your media content.

## Docker Compose

The easies way to use subtle is through docker. An example docker compose file is given below.

```yaml
---
services:
  subtle:
    image: nandeshdev/subtle:dev
    container_name: subtle
    environment:
      - SUBTLE_CONFIG_FILEPATH=/config/config.yaml
      - SUBTLE_DATABASE_FILEPATH=/config/database.db
      - SUBTLE_LOG_FILEPATH=/config/logs.log
      - SUBTLE_FILE_LOG_LEVEL=INFO
      - SUBTLE_CONSOLE_LOG_LEVEL=ERROR
    volumes:
      - /<config_directory_path>:/config
      - /<media_directory_path>:/media
    ports:
      - 3000:3000
    restart: unless-stopped
```

## License

[MIT](https://choosealicense.com/licenses/mit/)
