import { DirectoryView } from './directory_view'
import { VideoView } from './video_view'
import { useNavigation } from '@/src/utility/navigation'
import filepath from 'path-browserify'

export function Browser() {
    const navigation = useNavigation()
    const searchParams = navigation.useSearchParams()
    const pathSearchParam = searchParams.get('path')

    if (filepath.extname(pathSearchParam || '') == '') {
        return <DirectoryView />
    }

    return <VideoView />
}
