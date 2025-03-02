import { create } from '@bufbuild/protobuf'

import {
    CalculateSubtitleStatisticsRequest,
    CalculateSubtitleStatisticsResponseSchema,
    GetSubtitleCueOriginalDataRequest,
    GetSubtitleCueOriginalDataResponseSchema,
    GetSubtitleCueRequest,
    GetSubtitleCueResponseSchema,
    GetSubtitleCueSegmentRequest,
    GetSubtitleCueSegmentResponseSchema,
    GetSubtitleRequest,
    GetSubtitleResponseSchema,
    SubtitleOriginalFormat,
    SubtitleStage,
} from '../../../gen/proto/messages/subtitle_pb'

export async function calculateSubtitleStatistics({}: CalculateSubtitleStatisticsRequest) {
    return create(CalculateSubtitleStatisticsResponseSchema, {
        totalVideoCount: 100,
        videoWithExtractedSubtitleCount: 99,
        videoWithFormatedSubtitleCount: 96,
        videoWithExportedSubtitleCount: 68,
    })
}

export async function getSubtitle({ id }: GetSubtitleRequest) {
    return create(GetSubtitleResponseSchema, {
        id,
        title: 'Full Dialogue',
        stage: SubtitleStage.EXTRACTED,
        isProcessing: false,
        importIsExternal: false,
        originalFormat: SubtitleOriginalFormat.SRT,
        cueIds: ['0', '1', '2'],
    })
}

export async function getSubtitleCue({ id }: GetSubtitleCueRequest) {
    const index = parseInt(id)

    return create(GetSubtitleCueResponseSchema, {
        id,
        start: { seconds: BigInt(20 * index) },
        end: { seconds: BigInt(20 * index + 20) },
        segmentIds: ['0'],
    })
}

export async function getSubtitleCueOriginalData({}: GetSubtitleCueOriginalDataRequest) {
    return create(GetSubtitleCueOriginalDataResponseSchema)
}

export async function getSubtitleCueSegment({
    id,
}: GetSubtitleCueSegmentRequest) {
    return create(GetSubtitleCueSegmentResponseSchema, {
        id,
        text: 'Hey There!',
    })
}
