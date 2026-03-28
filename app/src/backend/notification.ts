import {notification} from "ant-design-vue";

notification.config({
    top: 96
})

class NotificationUtil {
    static success(message: string, description?: string) {
        notification.success({
            message,
            description,
        })
    }

    static error(message: string, description?: string) {
        notification.error({
            message,
            description,
        })
    }
}

export default NotificationUtil