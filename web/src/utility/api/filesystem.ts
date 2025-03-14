import { create } from '@bufbuild/protobuf'

import {
    GetRootDirectoryPathsRequest,
    GetRootDirectoryPathsResponseSchema,
    ReadDirectoryRequest,
    ReadDirectoryResponseSchema,
    SearchVideoRequest,
    SearchVideoResponseSchema,
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
            response.childrenDirectoryPaths = [
                '/media/series/Horimiya/Season 1',
                '/media/series/Horimiya/Season 2',
            ]
            response.videoPaths = [
                '/media/series/Horimiya/Horimiya - S00E00 - Special Bluray-1080p.mkv',
            ]
            break
        default:
            response.videoPaths = [
                '/media/series/Horimiya/Season 1/Horimiya - S01E01 - A Tiny Happenstance Bluray-1080p.mkv',
                '/media/series/Horimiya/Season 1/Horimiya - S01E02 - You Wear More Than One Face Bluray-1080p.mkv',
                "/media/series/Horimiya/Season 1/Horimiya - S01E03 - That's Why It's Okay Bluray-1080p.mkv",
                '/media/series/Horimiya/Season 1/Horimiya - S01E04 - Everybody Loves Somebody Bluray-1080p.mkv',
                "/media/series/Horimiya/Season 1/Horimiya - S01E05 - I Can't Say It Out Loud Bluray-1080p.mkv",
                "/media/series/Horimiya/Season 1/Horimiya - S01E06 - This Summer's Going to Be a Hot One Bluray-1080p.mkv",
                "/media/series/Horimiya/Season 1/Horimiya - S01E07 - You're Here, I'm Here Bluray-1080p.mkv",
                '/media/series/Horimiya/Season 1/Horimiya - S01E08 - The Truth Deception Reveals Bluray-1080p.mkv',
                "/media/series/Horimiya/Season 1/Horimiya - S01E09 - It's Hard, but Not Impossible Bluray-1080p.mkv",
                '/media/series/Horimiya/Season 1/Horimiya - S01E10 - Until the Snow Melts Bluray-1080p.mkv',
                '/media/series/Horimiya/Season 1/Horimiya - S01E11 - It May Seem Like Hate Bluray-1080p.mkv',
                '/media/series/Horimiya/Season 1/Horimiya - S01E12 - Hitherto, and Forevermore Bluray-1080p.mkv',
                '/media/series/Horimiya/Season 1/Horimiya - S01E13 - I Would Gift You the Sky Bluray-1080p.mkv',
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
