import { ProcessingIcon, TickIcon } from '../../../assets'
import { useProto } from '../../context/proto'
import { useQuery } from '@tanstack/react-query'
import { GetRoutinesRequest } from '../../../gen/proto/routine/routine_pb'

const REFETCH_INTERNAL = 10 * 1000

export function Routines() {
    const { RoutineServiceClient } = useProto()

    const { data, isLoading } = useQuery({
        queryKey: ['get-routines'],
        queryFn: () =>
            RoutineServiceClient?.getRoutines(new GetRoutinesRequest({})),
        refetchInterval: REFETCH_INTERNAL,
    })

    return (
        <section className="flex h-full flex-col gap-sm md:px-lg md:py-xxl">
            <div className="flex flex-row items-center gap-lg md:min-h-[4rem]">
                <h2 className="text-md text-gray-830">Routines</h2>
            </div>
            <p className="text-sm text-gray-830">
                Routines are small tasks which run in background. Specific user
                actions will be disabled when a routine is running.
            </p>
            <div className="h-[4px] w-full rounded-sm bg-gray-80" />
            <section className="flex flex-col gap-sm overflow-y-auto">
                {isLoading ? (
                    <div className="flex grid-cols-[2fr_1fr] flex-col gap-xs rounded-sm bg-gray-80 p-sm">
                        <div className="flex flex-row justify-between">
                            <div className="h-md w-[16rem] animate-pulse rounded-sm bg-gray-190" />
                            <div className="h-md w-[4rem] animate-pulse rounded-sm bg-gray-120" />
                        </div>
                        <div className="h-sm w-4/5 animate-pulse rounded-sm bg-gray-120" />
                        <div className="h-sm w-1/2 animate-pulse rounded-sm bg-gray-120" />
                    </div>
                ) : (
                    data?.routines.map((routine) => {
                        return (
                            <div
                                className="flex grid-cols-[2fr_1fr] flex-col gap-xs rounded-sm bg-gray-80 p-sm"
                                key={routine.name}
                            >
                                <div className="flex flex-row justify-between">
                                    <p className="text-start text-sm text-gray-830">
                                        {routine.name}
                                    </p>
                                    {routine.isRunning ? (
                                        <ProcessingIcon />
                                    ) : (
                                        <TickIcon />
                                    )}
                                </div>
                                <p className="text-sm text-gray-520">
                                    {routine.description}
                                </p>
                            </div>
                        )
                    })
                )}
            </section>
        </section>
    )
}
