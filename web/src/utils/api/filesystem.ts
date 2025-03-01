import { create } from '@bufbuild/protobuf'
import {
    GetRootDirectoryPathsRequest,
    GetRootDirectoryPathsResponseSchema,
    SearchVideoRequest,
    SearchVideoResponseSchema,
    ReadDirectoryRequest,
    ReadDirectoryResponseSchema,
} from '../../../gen/proto/messages/filesystem_pb'

export async function getRootDirectoryPaths({}: GetRootDirectoryPathsRequest) {
    return create(GetRootDirectoryPathsResponseSchema, {
        paths: ['/media/series'],
    })
}

export async function readDirectory({ path }: ReadDirectoryRequest) {
    const response = create(ReadDirectoryResponseSchema)

    switch (path) {
        case '/media/series':
            response.childrenDirectoryPaths = ['/media/series/Horimiya']
            break
        case '/media/series/Horimiya':
            response.childrenDirectoryPaths = ['/media/series/Horimiya/Season 1']
            break
        default:
            response.videoPaths = [
                '/media/series/Horimiya/Season 1/Horimiya - S01E01 - A Tiny Happenstance Bluray-1080p.mkv',
                '/media/series/Horimiya/Season 1/Horimiya - S01E02 - You Wear More Than One Face Bluray-1080p.mkv',
                "/media/series/Horimiya/Season 1/Horimiya - S01E03 - That's Why It's Okay Bluray-1080p.mkv",
                '/media/series/Horimiya/Season 1/Horimiya - S01E04 - Everybody Loves Somebody Bluray-1080p.mkv',
            ]
    }

    return response
}

export async function searchVideo({ path }: SearchVideoRequest) {
    return create(SearchVideoResponseSchema, {
        id: '0',
        path,
        subtitleIds: ['0', '1', '2'],
    })
}
