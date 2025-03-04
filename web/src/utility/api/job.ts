import { create } from '@bufbuild/protobuf'
import { timestampNow } from '@bufbuild/protobuf/wkt'

import {
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
