<template>
  <ContentBase>
    <el-container v-if="dataLoaded">
      <el-main>
        <el-table :data="books">
          <el-table-column label="名称" prop="title"/>
          <el-table-column label="分类" prop="category"/>
          <el-table-column label="创建时间" prop="created_at" :formatter="formatDateTime"/>
          <el-table-column label="更新时间" prop="updated_at" :formatter="formatDateTime"/>
          <el-table-column label="操作">
            <template #default="scope">
              <el-button @click="handleViewEntries(scope.row)">查看词汇</el-button>
              <el-button @click="handleEditBook(scope.row)">编辑</el-button>
              <el-button @click="handleDeleteBook(scope.row)">删除</el-button>
            </template>
          </el-table-column>
        </el-table>
        <div>
          <el-pagination
              v-model:current-page="currentPage"
              v-model:page-size="pageSize"
              :page-sizes="[5, 10, 20, 50, 100]"
              :background="true"
              layout="total, sizes, prev, pager, next, jumper"
              :total="countTotalRecords"
              @size-change="handleSizeChange"
              @current-change="handleCurrentChange"
          />
        </div>
      </el-main>
      <br>
      <el-footer>
        <el-button type="primary" @click="showAddDialog">添加</el-button>
      </el-footer>
    </el-container>

    <el-dialog v-model="editFormVisible" title="编辑词书">
      <el-form :model="book">
        <el-form-item label="标题">
          <el-input v-model="book.title" @keyup.enter="confirmUpdate" autocomplete="off"/>
        </el-form-item>
        <el-form-item label="分类">
          <el-input v-model="book.category" @keyup.enter="confirmUpdate" autocomplete="off"/>
        </el-form-item>
      </el-form>
      <template #footer>
                <span class="dialog-footer">
                  <el-button @click="cancelUpdate">取消</el-button>
                  <el-button type="primary" @click="confirmUpdate">确认</el-button>
                </span>
      </template>
    </el-dialog>

    <el-dialog v-model="addFormVisible" title="新建词书">
      <el-form :model="newBook">
        <el-form-item label="标题">
          <el-input v-model="newBook.title" autocomplete="off"/>
        </el-form-item>
        <el-form-item label="分类">
          <el-input v-model="newBook.category" @keyup.enter="confirmAdd" autocomplete="off"/>
        </el-form-item>
        <el-checkbox v-model="uploadEnabled">上传文件(如不上传文件，则创建一个空的词书)</el-checkbox>
        <el-form-item v-if="uploadEnabled">
          <el-upload
              ref="upload"
              :before-upload="beforeUpload"
              :action="uploadFileAPI"
              :headers="headers"
              :data="{title: newBook.title, category: newBook.category}"
              :name="`file`"
              :limit="1"
              :on-exceed="handleExceed"
              :auto-upload="false"
              :on-success="handleAddBookSuccess"
          >
            <template #trigger>
              <el-button type="primary">选择文件</el-button>
            </template>
            <template #tip>
              <div style="color: red">
                一次只能上传一个.txt文件或者.xlsx文件，词条数不能超过50万个
              </div>
            </template>
          </el-upload>
        </el-form-item>
      </el-form>
      <template #footer>
                <span class="dialog-footer">
                  <el-button @click="cancelAdd">Cancel</el-button>
                  <el-button type="primary" @click="confirmAdd">确认</el-button>
                </span>
      </template>
    </el-dialog>
  </ContentBase>
</template>

<script setup>
import {
  ElButton,
  ElCheckbox,
  ElContainer,
  ElDialog,
  ElFooter,
  ElForm,
  ElFormItem,
  ElInput,
  ElMain,
  ElNotification,
  ElPagination,
  ElTable,
  ElTableColumn,
  ElUpload,
  genFileId,
} from 'element-plus';
import "element-plus/dist/index.css";
import {useStore} from "vuex";
import ContentBase from "@/components/ContentBase.vue";
import {onMounted, reactive, ref} from "vue";
import {getBookCount, getBookList} from "@/assets/js/module/book/query";
import {checkBookFileType} from "@/assets/js/util/file_util";
import {updateBook} from "@/assets/js/module/book/update";
import {deleteBook} from "@/assets/js/module/book/delete";
import {formatDateTime} from "@/assets/js/util/datetime_util";
import {addBook} from "@/assets/js/module/book/add";
import router from "@/router";

const store = useStore();
const checkResponse = (response) => {
  if (response == null || response.code === 2) {
    store.dispatch("logout")
    location.reload();
  }
}
let editFormVisible = ref(false);

let book = ref(null)

let books = reactive([]);
let pageSize = ref(5);
let currentPage = ref(1);

const handleCurrentChange = async () => {
  dataLoaded.value = false;
  const getBookListResponse = await getBookList(pageSize.value, currentPage.value);
  checkResponse(getBookListResponse)
  books = getBookListResponse.data;
  dataLoaded.value = true;
};

const handleSizeChange = async () => {
  dataLoaded.value = false;
  const getBookListResponse = await getBookList(pageSize.value, currentPage.value);
  checkResponse(getBookListResponse)
  books = getBookListResponse.data;
  dataLoaded.value = true;
}

const dataLoaded = ref(false);
let countTotalRecords = ref(0);
onMounted(
    async () => {
      const getBookCountResponse = await getBookCount();
      checkResponse(getBookCountResponse);
      countTotalRecords.value = getBookCountResponse.data;

      const getBookListResponse = await getBookList(pageSize.value, currentPage.value);
      checkResponse(getBookListResponse);
      books = getBookListResponse.data;
      dataLoaded.value = true;
    }
);

const handleViewEntries = (row) => {
  localStorage.setItem("book_id", row.id);
  router.push({name: "word_setting"});
}
const handleEditBook = (row) => {
  book.value = row;
  editFormVisible.value = true;
}


const cancelUpdate = () => {
  editFormVisible.value = false;
  location.reload();
}
const confirmUpdate = async () => {
  editFormVisible.value = false;
  dataLoaded.value = false;
  const updateBookResponse = await updateBook(book.value.id, book.value.title, book.value.category, book.value.created_at);
  checkResponse(updateBookResponse);
  dataLoaded.value = true;
  location.reload();
}
const handleDeleteBook = async (row) => {
  ElNotification({
    title: '删除词书',
    message: '词书删除中, 请稍等...',
    type: 'info',
    duration: 1000,
  });
  const deleteBookResponse = await deleteBook(row.id);
  checkResponse(deleteBookResponse);
  if (deleteBookResponse.code === 0) {
    dataLoaded.value = false;
    const getBookCountResponse = await getBookCount();
    checkResponse(getBookCountResponse);
    countTotalRecords.value = getBookCountResponse.data;
    const getBookListResponse = await getBookList(pageSize.value, currentPage.value);
    checkResponse(getBookListResponse);
    books = getBookListResponse.data;
    dataLoaded.value = true;
    ElNotification({
      title: 'Success',
      message: '删除成功',
      type: 'success',
      duration: 1000,
    })
  }
}


const addFormVisible = ref(false);
const newBook = reactive({
  title: '',
  category: '',
});

const showAddDialog = () => {
  addFormVisible.value = true;
}

const uploadEnabled = ref(false);
const upload = ref(null);

const handleExceed = (files) => {
  upload.value.clearFiles()
  const file = files[0]
  file.uid = genFileId()
  upload.value.handleStart(file)
}

const handleAddBookSuccess = (response) => {
  if (response.code === 1) {
    ElNotification({
      title: '上传失败',
      message: response.message,
      type: 'error',
      duration: 1000,
    })
  } else {
    ElNotification({
      title: 'Success',
      message: response.message,
      type: 'success',
      duration: 1000,
    })
  }
  setTimeout(() => {
    location.reload();
  }, 1000);
}

const handleAddBookError = (response) => {
  ElNotification({
    title: 'Error',
    message: response.message,
    type: 'error',
    duration: 1000,
  })
}

const cancelAdd = () => {
  uploadEnabled.value = false;
  addFormVisible.value = false;
  newBook.title = '';
  newBook.category = '';
}


const confirmAdd = async () => {
  ElNotification({
    title: '文件上传',
    message: '词书上传中, 请稍等...',
    type: 'info',
    duration: 1000,
  });
  addFormVisible.value = false;
  uploadEnabled.value = false;
  dataLoaded.value = false;
  if (upload.value != null) {
    upload.value.submit();
  } else {
    const addBookResponse = await addBook(newBook.title, newBook.category);
    checkResponse(addBookResponse);
    if (addBookResponse.code === 0) {
      ElNotification({
        title: 'Success',
        message: '添加成功',
        type: 'success',
        duration: 1000,
      });
    } else {
      handleAddBookError(addBookResponse);
    }
  }
  const getBookCountResponse = await getBookCount();
  checkResponse(getBookCountResponse);
  countTotalRecords.value = getBookCountResponse.data.data;
  const getBookListResponse = await getBookList(pageSize.value, currentPage.value);
  checkResponse(getBookListResponse);
  books = getBookListResponse.data;
  dataLoaded.value = true;
  newBook.title = '';
  newBook.category = '';
}

const serverLink = localStorage.getItem('server_link');
const uploadFileAPI = ref(serverLink + '/book/add');
const headers = ref({
  'Authorization': 'Bearer ' + localStorage.getItem('token')
})

const beforeUpload = (file) => {
  const res = checkBookFileType(file);
  if (res.code === 0) {
    return true;
  } else {
    ElNotification({
      title: 'Error',
      message: res.message,
      type: 'error',
      duration: 1000,
    })
    return false;
  }
}
</script>

<style scoped>
.dialog-footer button:first-child {
  margin-right: 10px;
}
</style>