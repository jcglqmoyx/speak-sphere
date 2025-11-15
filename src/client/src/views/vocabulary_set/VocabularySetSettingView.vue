<template>
  <ContentBase>
    <el-container v-if="dataLoaded">
      <el-main>
        <el-table :data="vocabularySets">
          <el-table-column label="名称" prop="title"/>
          <el-table-column label="分类" prop="category"/>
          <el-table-column label="创建时间" prop="created_at" :formatter="formatDateTime"/>
          <el-table-column label="更新时间" prop="updated_at" :formatter="formatDateTime"/>
          <el-table-column label="操作">
            <template #default="scope">
              <el-button @click="handleViewVocabularies(scope.row)">查看词汇</el-button>
              <el-button @click="handleEditVocabularySet(scope.row)">编辑</el-button>
              <el-button @click="handleDeleteVocabularySet(scope.row)">删除</el-button>
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
      <el-form :model="vocabularySet">
        <el-form-item label="标题">
          <el-input v-model="vocabularySet.title" @keyup.enter="confirmUpdate" autocomplete="off"/>
        </el-form-item>
        <el-form-item label="分类">
          <el-input v-model="vocabularySet.category" @keyup.enter="confirmUpdate" autocomplete="off"/>
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
      <el-form :model="newVocabularySet">
        <el-form-item label="标题">
          <el-input v-model="newVocabularySet.title" autocomplete="off"/>
        </el-form-item>
        <el-form-item label="分类">
          <el-input v-model="newVocabularySet.category" @keyup.enter="confirmAdd" autocomplete="off"/>
        </el-form-item>
        <el-checkbox v-model="uploadEnabled">上传文件(如不上传文件，则创建一个空的词书)</el-checkbox>
        <el-form-item v-if="uploadEnabled">
          <el-upload
              ref="upload"
              :before-upload="beforeUpload"
              :action="uploadFileAPI"
              :headers="headers"
              :data="{title: newVocabularySet.title, category: newVocabularySet.category}"
              :name="`file`"
              :limit="1"
              :on-exceed="handleExceed"
              :auto-upload="false"
              :on-success="handleAddVocabularySetSuccess"
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
import {getVocabularySetCount, getVocabularySetList} from "@/assets/js/module/vocabulary_set/query";
import {checkVocabularySetFileType} from "@/assets/js/util/file_util";
import {updateVocabularySet} from "@/assets/js/module/vocabulary_set/update";
import {deleteVocabularySet} from "@/assets/js/module/vocabulary_set/delete";
import {formatDateTime} from "@/assets/js/util/datetime_util";
import {addVocabularySet} from "@/assets/js/module/vocabulary_set/add";
import router from "@/router";

const store = useStore();
const checkResponse = (response) => {
  if (response == null || response.code === 2) {
    store.dispatch("logout")
    location.reload();
  }
}
let editFormVisible = ref(false);

let vocabularySet = ref(null)

let vocabularySets = reactive([]);
let pageSize = ref(5);
let currentPage = ref(1);

const handleCurrentChange = async () => {
  dataLoaded.value = false;
  const getVocabularySetListResponse = await getVocabularySetList(pageSize.value, currentPage.value);
  checkResponse(getVocabularySetListResponse)
  vocabularySets = getVocabularySetListResponse.data;
  dataLoaded.value = true;
};

const handleSizeChange = async () => {
  dataLoaded.value = false;
  const getVocabularySetListResponse = await getVocabularySetList(pageSize.value, currentPage.value);
  checkResponse(getVocabularySetListResponse)
  vocabularySets = getVocabularySetListResponse.data;
  dataLoaded.value = true;
}

const dataLoaded = ref(false);
let countTotalRecords = ref(0);
onMounted(
    async () => {
      const getVocabularySetCountResponse = await getVocabularySetCount();
      checkResponse(getVocabularySetCountResponse);
      countTotalRecords.value = getVocabularySetCountResponse.data;

      const getVocabularySetListResponse = await getVocabularySetList(pageSize.value, currentPage.value);
      checkResponse(getVocabularySetListResponse);
      vocabularySets = getVocabularySetListResponse.data;
      dataLoaded.value = true;
    }
);

const handleViewVocabularies = (row) => {
  localStorage.setItem("vocabulary_set_id", row.id);
  router.push({name: "vocabulary_setting"});
}
const handleEditVocabularySet = (row) => {
  vocabularySet.value = row;
  editFormVisible.value = true;
}


const cancelUpdate = () => {
  editFormVisible.value = false;
  location.reload();
}
const confirmUpdate = async () => {
  editFormVisible.value = false;
  dataLoaded.value = false;
  const updateVocabularySetResponse = await updateVocabularySet(vocabularySet.value.id, vocabularySet.value.title, vocabularySet.value.category, vocabularySet.value.created_at);
  checkResponse(updateVocabularySetResponse);
  dataLoaded.value = true;
  location.reload();
}
const handleDeleteVocabularySet = async (row) => {
  ElNotification({
    title: '删除词书',
    message: '词书删除中, 请稍等...',
    type: 'info',
    duration: 1000,
  });
  const deleteVocabularySetResponse = await deleteVocabularySet(row.id);
  checkResponse(deleteVocabularySetResponse);
  if (deleteVocabularySetResponse.code === 0) {
    dataLoaded.value = false;
    const getVocabularySetCountResponse = await getVocabularySetCount();
    checkResponse(getVocabularySetCountResponse);
    countTotalRecords.value = getVocabularySetCountResponse.data;
    const getVocabularySetListResponse = await getVocabularySetList(pageSize.value, currentPage.value);
    checkResponse(getVocabularySetListResponse);
    vocabularySets = getVocabularySetListResponse.data;
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
const newVocabularySet = reactive({
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

const handleAddVocabularySetSuccess = (response) => {
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

const handleAddVocabularySetError = (response) => {
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
  newVocabularySet.title = '';
  newVocabularySet.category = '';
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
    const addVocabularySetResponse = await addVocabularySet(newVocabularySet.title, newVocabularySet.category);
    checkResponse(addVocabularySetResponse);
    if (addVocabularySetResponse.code === 0) {
      ElNotification({
        title: 'Success',
        message: '添加成功',
        type: 'success',
        duration: 1000,
      });
    } else {
      handleAddVocabularySetError(addVocabularySetResponse);
    }
  }
  const getVocabularySetCountResponse = await getVocabularySetCount();
  checkResponse(getVocabularySetCountResponse);
  countTotalRecords.value = getVocabularySetCountResponse.data.data;
  const getVocabularySetListResponse = await getVocabularySetList(pageSize.value, currentPage.value);
  checkResponse(getVocabularySetListResponse);
  vocabularySets = getVocabularySetListResponse.data;
  dataLoaded.value = true;
  newVocabularySet.title = '';
  newVocabularySet.category = '';
}

const serverLink = localStorage.getItem('server_link');
const uploadFileAPI = ref(serverLink + '/vocabulary_set/add');
const headers = ref({
  'Authorization': 'Bearer ' + localStorage.getItem('token')
})

const beforeUpload = (file) => {
  const res = checkVocabularySetFileType(file);
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