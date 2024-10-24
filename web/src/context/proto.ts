import { Client } from '@connectrpc/connect'
import { MediaService } from '../../gen/proto/media/media_connect'
import { createContext, useContext } from 'react'

export type ProtoContent = {
    MediaServiceClient?: Client<typeof MediaService>
}

export const ProtoContext = createContext<ProtoContent>({})

export const useProto = () => {
    return useContext(ProtoContext)
}
