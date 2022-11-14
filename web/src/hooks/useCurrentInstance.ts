export default function useCurrentInstance() {
    // const { appContext } = getCurrentInstance() as ComponentInternalInstance
    // const proxy = appContext.config.globalProperties
    const { proxy } = getCurrentInstance() as ComponentInternalInstance
    return {
        proxy
    }
}
