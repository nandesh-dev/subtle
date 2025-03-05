import { create } from '@bufbuild/protobuf'
import { timestampNow } from '@bufbuild/protobuf/wkt'

import {
    GetJobLogRequest,
    GetJobLogResponseSchema,
    GetJobLogsRequest,
    GetJobLogsResponseSchema,
    GetJobRequest,
    GetJobResponseSchema,
    GetJobsRequest,
    GetJobsResponseSchema,
} from '@/gen/proto/messages/job_pb'

export async function getJobs({}: GetJobsRequest) {
    return create(GetJobsResponseSchema, {
        jobIds: ['scan', 'extract', 'format', 'export'],
    })
}

export async function getJob({ id }: GetJobRequest) {
    switch (id) {
        case 'extract':
            return create(GetJobResponseSchema, {
                id,
                sequenceNumber: 2,
                name: 'Extracting',
                description:
                    'A really long description about the job which tells user what it does.',
                isRunning: false,
                lastRun: timestampNow(),
            })
        case 'format':
            return create(GetJobResponseSchema, {
                id,
                sequenceNumber: 3,
                name: 'Formating',
                description:
                    'A really long description about the job which tells user what it does.',
                isRunning: true,
                lastRun: timestampNow(),
            })
        case 'export':
            return create(GetJobResponseSchema, {
                id,
                sequenceNumber: 4,
                name: 'Exporting',
                description:
                    'A really long description about the job which tells user what it does.',
                isRunning: false,
                lastRun: timestampNow(),
            })
        default:
            return create(GetJobResponseSchema, {
                id,
                sequenceNumber: 1,
                name: 'Scanning',
                description:
                    'A really long description about the job which tells user what it does.',
                isRunning: false,
                lastRun: timestampNow(),
            })
    }
}

export async function getJobLogs({
    limit,
    olderThanLogId,
    newerThanLogId,
}: GetJobLogsRequest) {
    const STARTING_LOG_ID = 100
    const MAXIMUM_LOG_ID = 110

    const ids = []

    let start = STARTING_LOG_ID
    let end = 0

    if (
        (olderThanLogId == undefined || olderThanLogId == '') &&
        (newerThanLogId == undefined || newerThanLogId == '') &&
        limit !== undefined
    ) {
        end = STARTING_LOG_ID - limit + 1
    } else if (
        olderThanLogId !== undefined &&
        olderThanLogId !== '' &&
        newerThanLogId !== undefined &&
        newerThanLogId !== ''
    ) {
        start = parseInt(olderThanLogId) - 1
        end = parseInt(newerThanLogId) + 1
    } else if (olderThanLogId !== undefined && newerThanLogId !== '') {
        start = parseInt(olderThanLogId) - 1

        if (limit !== undefined) {
            end = parseInt(olderThanLogId) - limit + 1
        }
    } else if (newerThanLogId !== undefined && newerThanLogId !== '') {
        end = parseInt(newerThanLogId) + 1

        if (limit == undefined) {
            start = MAXIMUM_LOG_ID
        } else {
            start = Math.min(
                parseInt(newerThanLogId) + limit - 1,
                MAXIMUM_LOG_ID
            )
        }
    }

    for (let i = start; i >= end; i--) {
        ids.push(i.toString())
    }

    return create(GetJobLogsResponseSchema, { ids })
}

export async function getJobLog({ id }: GetJobLogRequest) {
    if (id == '0') {
        return create(GetJobLogResponseSchema, {
            id,
            jobId: 'scan',
            jobName: 'Scanning',
            isSuccess: false,
            errorMessage: 'Failed to read file',
            startTimestamp: timestampNow(),
            duration: { seconds: BigInt(10000) },
        })
    }

    return create(GetJobLogResponseSchema, {
        id,
        jobId: 'scan',
        jobName: 'Scanning',
        isSuccess: true,
        startTimestamp: timestampNow(),
        duration: { seconds: BigInt(10000) },
    })
}
