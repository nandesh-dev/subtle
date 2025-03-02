import { Interceptor } from '@connectrpc/connect'

async function mockApiDelay() {
  await fetch("/")
}

export const mockApiDelayInterceptor: Interceptor = (next) => async (req) => {
    await mockApiDelay()
    return await next(req)
}
