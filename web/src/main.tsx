import { Routes, useNavigation } from './utils/navigation'

export function Main() {
    const navigation = useNavigation()
    const route = navigation?.useRoute()

    return (
        <div>
            {route}
            <button onClick={()=>{
              navigation?.navigate(Routes.Jobs)
            }}>Update</button>
        </div>
    )
}
