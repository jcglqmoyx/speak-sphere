import {dayjs} from "element-plus";

const formatDate = (row) => {
    const date = row['date_to_review'].toString();
    return date.substring(0, 4) + "-" + date.substring(4, 6) + "-" + date.substring(6);
}
const formatDateTime = (row, column) => {
    return dayjs(row[column.property]).format('YYYY-MM-DD HH:mm:ss');
}

export {formatDate, formatDateTime}