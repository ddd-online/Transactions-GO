import { message, notification } from "ant-design-vue";

notification.config({
    top: 96
})

class NotificationUtil {
    static success(messageText: string, description?: string) {
        if (description) {
            notification.success({
                message: messageText,
                description,
            })
        } else {
            message.success(messageText);
        }
    }

    static error(messageText: string, description?: string) {
        if (description) {
            notification.error({
                message: messageText,
                description,
            })
        } else {
            message.error(messageText);
        }
    }
}

export default NotificationUtil
