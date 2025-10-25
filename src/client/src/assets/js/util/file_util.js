export function checkBookFileType(file) {
    if (file.type === 'text/plain' || file.type === 'application/vnd.openxmlformats-officedocument.spreadsheetml.sheet') {
        if (file.size > 1024 * 1024 * 10) {
            return {
                code: 1,
                message: "文件大小不能超过10M",
            };
        }
        return {
            code: 0,
            message: "success",
        };
    } else {
        return {
            code: 1,
            message: "文件格式错误",
        };
    }
}