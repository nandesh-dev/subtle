import { create } from '@bufbuild/protobuf'
import {
    GetConfigRequest,
    GetConfigResponseSchema,
    UpdateConfigRequest,
    UpdateConfigResponseSchema,
} from '../../../gen/proto/messages/config_pb'

export async function getConfig({}: GetConfigRequest) {
    return create(GetConfigResponseSchema, {
        config: `watch_directories:
    - path: /media
      extraction:
        enable: true
        raw_stream_title_keywords:
            - Full
            - Dialogue
        formats:
            ass:
                enable: true
                languages:
                    - en
            pgs:
                enable: true
                languages:
                    - en
      formating:
        enable: true
        text_based_subtitle:
            charactor_mappings: []
        image_based_subtitle:
            charactor_mappings:
                - language: en
                  mappings:
                    - from: '|'
                      to: I
      exporting:
        enable: true
        format: srt
server:
    web:
        port: 3000
        cor_origins: []
        enable_grpc_reflection: false
        serve_directory: /public
    database:
        path: testing/config/database.db
    job:
        delay: 15m0s
    logging:
        path: testing/config/logs.log
        console_level: DEBUG
        file_level: INFO
`,
    })
}

export async function updateConfig({}:UpdateConfigRequest){
  return create(UpdateConfigResponseSchema)
}
