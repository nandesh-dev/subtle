import { Client } from '@connectrpc/connect'
import { MediaService } from '../../gen/proto/media/media_connect'
import { createContext, useContext } from 'react'
import { SubtitleService } from '../../gen/proto/subtitle/subtitle_connect'

export type ProtoContent = {
    MediaServiceClient?: Client<typeof MediaService>
    SubtitleServiceClient?: Client<typeof SubtitleService>
}

export const ProtoContext = createContext<ProtoContent>({})

export const useProto = () => {
    return useContext(ProtoContext)
}
