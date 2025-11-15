<template>
  <ContentBase>
    <el-container v-if="dataLoaded">
      <el-main>
        <el-table :data="vocabularies">
          <el-table-column label="单词" prop="vocabulary"/>
          <el-table-column label="含义" prop="meaning"/>
          <el-table-column label="备注" prop="note"/>
          <el-table-column label="不想学" prop="unwanted"/>
          <el-table-column label="学习次数" prop="study_count"/>
          <el-table-column label="复习日期" prop="date_to_review" :formatter="formatDate"/>
          <el-table-column label="创建时间" prop="created_at" :formatter="formatDateTime"/>
          <el-table-column label="更新时间" prop="updated_at" :formatter="formatDateTime"/>
          <el-table-column label="操作">
            <template #default="scope">
              <el-tooltip content="编辑单词条目" placement="top">
                <el-button @click="handleEditVocabulary(scope.row)">编辑</el-button>
              </el-tooltip>
              <el-tooltip content="删除单词条目" placement="top">
                <el-button @click="handleDeleteVocabulary(scope.row)">删除</el-button>
              </el-tooltip>
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
        <el-tooltip content="添加新单词" placement="top">
          <el-button type="primary" @click="showAddDialog">添加</el-button>
        </el-tooltip>
      </el-footer>
    </el-container>


    <el-dialog v-model="editFormVisible" title="编辑条目">
      <el-form :model="vocabulary">
        <el-form-item label="单词">
          <el-input v-model="vocabulary.vocabulary" autocomplete="off"/>
        </el-form-item>
        <el-form-item label="含义">
          <el-input type="textarea" :autosize="{minRows: 5}" v-model="vocabulary.meaning" autocomplete="off"/>
        </el-form-item>
        <el-form-item label="备注">
          <el-input type="textarea" :autosize="{minRows: 3}" v-model="vocabulary.note" autocomplete="off"/>
        </el-form-item>
        <el-form-item>
          <el-checkbox v-model="vocabulary.unwanted" label="标记为不想学" size="large"/>
        </el-form-item>
        <el-form-item label="学习次数">
          <el-input-number v-model="vocabulary.study_count" :min="0" :max="30"/>
        </el-form-item>
        <el-form-item label="下次复习日期">
          <div class="block">
            <el-date-picker
                v-model="dateToReviewEdit"
                :disabled-date="disabledDateToReview"
                type="date"
                placeholder="选择日期"
            />
          </div>
        </el-form-item>
      </el-form>
      <template #footer>
                <span class="dialog-footer">
                  <el-tooltip content="取消编辑" placement="top">
                    <el-button @click="cancelUpdate">取消</el-button>
                  </el-tooltip>
                  <el-tooltip content="确认编辑" placement="top">
                    <el-button type="primary" @click="confirmUpdate">确认</el-button>
                  </el-tooltip>
                </span>
      </template>
    </el-dialog>

    <el-dialog v-model="addFormVisible" title="新建条目">
      <el-form :model="newVocabulary">
        <el-form-item label="单词">
          <el-input v-model="newVocabulary.vocabulary" @keyup.enter="confirmAdd" autocomplete="off"/>
        </el-form-item>
        <el-form-item label="含义">
          <el-input v-model="newVocabulary.meaning" @keyup.enter="confirmAdd" autocomplete="off"/>
        </el-form-item>
        <el-form-item label="备注">
          <el-input v-model="newVocabulary.note" @keyup.enter="confirmAdd" autocomplete="off"/>
        </el-form-item>
      </el-form>
      <template #footer>
                <span class="dialog-footer">
                  <el-tooltip content="取消添加" placement="top">
                    <el-button @click="cancelAdd">Cancel</el-button>
                  </el-tooltip>
                  <el-tooltip content="确认添加 (快捷键: Enter)" placement="top">
                    <el-button type="primary" @click="confirmAdd">确认</el-button>
                  </el-tooltip>
                </span>
      </template>
    </el-dialog>
  </ContentBase>
</template>

<script setup>
import "element-plus/dist/index.css"
import {formatDate, formatDateTime} from "@/assets/js/util/datetime_util";
import {
  ElButton,
  ElCheckbox,
  ElContainer,
  ElDatePicker,
  ElDialog,
  ElFooter,
  ElForm,
  ElFormItem,
  ElInput,
  ElInputNumber,
  ElMain,
  ElNotification,
  ElPagination,
  ElTable,
  ElTableColumn,
  ElTooltip,
} from "element-plus";
import {onMounted, reactive, ref} from "vue";
import ContentBase from "@/components/ContentBase.vue";
import {getVocabularyCount, getVocabularyList} from "@/assets/js/module/vocabulary/query";
import {useStore} from "vuex";
import {updateVocabulary} from "@/assets/js/module/vocabulary/update";
import {deleteVocabulary} from "@/assets/js/module/vocabulary/delete";
import {AddVocabulary} from "@/assets/js/module/vocabulary/add";

const vocabularySetId = parseInt(localStorage.getItem("vocabulary_set_id"));
let vocabularies = reactive([]);
let pageSize = ref(5);
let currentPage = ref(1);

const store = useStore();
const checkResponse = (response) => {
  if (response == null || response.code === 2) {
    store.dispatch("logout")
    location.reload();
  }
}

const dataLoaded = ref(false);
let countTotalRecords = ref(0);
onMounted(
    async () => {
      const getVocabularyCountResponse = await getVocabularyCount(vocabularySetId);
      checkResponse(getVocabularyCountResponse);
      countTotalRecords.value = getVocabularyCountResponse.data;

      const getVocabularyListResponse = await getVocabularyList(vocabularySetId, pageSize.value, currentPage.value);
      checkResponse(getVocabularyListResponse);
      vocabularies = getVocabularyListResponse.data;
      dataLoaded.value = true;
    }
);

const handleCurrentChange = async () => {
  dataLoaded.value = false;
  const getVocabularyListResponse = await getVocabularyList(vocabularySetId, pageSize.value, currentPage.value);
  checkResponse(getVocabularyListResponse);
  vocabularies = getVocabularyListResponse.data;
  dataLoaded.value = true;
};

const handleSizeChange = async () => {
  dataLoaded.value = false;
  const getVocabularyListResponse = await getVocabularyList(vocabularySetId, pageSize.value, currentPage.value);
  checkResponse(getVocabularyListResponse)
  vocabularies = getVocabularyListResponse.data;
  dataLoaded.value = true;
}

const vocabulary = ref(null);
const editFormVisible = ref(false);


const handleDeleteVocabulary = async (row) => {
  ElNotification({
    title: '删除条目',
    message: '条目删除中, 请稍等...',
    type: 'info',
    duration: 1000,
  });
  const deleteVocabularyResponse = await deleteVocabulary(row.id);
  checkResponse(deleteVocabularyResponse);
  if (deleteVocabularyResponse.code === 0) {
    dataLoaded.value = false;
    const getVocabularyCountResponse = await getVocabularyCount(vocabularySetId);
    checkResponse(getVocabularyCountResponse);
    countTotalRecords.value = getVocabularyCountResponse.data;
    const getVocabularyListResponse = await getVocabularyList(vocabularySetId, pageSize.value, currentPage.value);
    checkResponse(getVocabularyListResponse);
    vocabularies = getVocabularyListResponse.data;
    dataLoaded.value = true;
    ElNotification({
      title: 'Success',
      message: '删除成功',
      type: 'success',
      duration: 1000,
    });
  }
}

const disabledDateToReview = (date) => {
  const year = date.getFullYear();
  return year < 2000 || year > 2099;
}

const dateToReviewEdit = ref(0)
const handleEditVocabulary = (row) => {
  vocabulary.value = row;
  const dateNumber = row.date_to_review;
  const year = Math.floor(dateNumber / 10000);
  const month = Math.floor((dateNumber % 10000) / 100) - 1; // 月份是从0开始的，所以需要减1
  const day = dateNumber % 100;
  const dateObj = new Date(year, month, day);
  dateToReviewEdit.value = dateObj.getTime();
  editFormVisible.value = true;
}

const cancelUpdate = () => {
  editFormVisible.value = false;
  location.reload();
}

const confirmUpdate = async () => {
  editFormVisible.value = false;
  dataLoaded.value = false;

  const date = new Date(dateToReviewEdit.value);
  const year = date.getFullYear();
  const month = date.getMonth() + 1;
  const day = date.getDate();

  vocabulary.value.date_to_review = year * 10000 + month * 100 + day;

  const updateVocabularyResponse = await updateVocabulary(
      vocabulary.value.id,
      vocabulary.value.vocabulary,
      vocabulary.value.meaning,
      vocabulary.value.vocabulary_set_id,
      vocabulary.value.note,
      vocabulary.value.unwanted,
      vocabulary.value.study_count,
      vocabulary.value.date_to_review,
      vocabulary.value.created_at,
  )
  checkResponse(updateVocabularyResponse);
  dataLoaded.value = true;
}

const addFormVisible = ref(false);
const showAddDialog = () => {
  addFormVisible.value = true;
}

const newVocabulary = reactive({
  vocabulary: '',
  meaning: '',
  note: '',
})
const cancelAdd = () => {
  addFormVisible.value = false;
  newVocabulary.vocabulary = '';
  newVocabulary.meaning = '';
  newVocabulary.note = '';
}
const handleAddVocabularyError = (response) => {
  ElNotification({
    title: 'Error',
    message: response.message,
    type: 'error',
    duration: 1000,
  })
}
const confirmAdd = async () => {
  addFormVisible.value = false;
  dataLoaded.value = false;
  const addVocabularyResponse = await AddVocabulary(vocabularySetId, newVocabulary.vocabulary, newVocabulary.meaning, newVocabulary.note);
  checkResponse(addVocabularyResponse);
  if (addVocabularyResponse.code === 0) {
    ElNotification({
      title: 'Success',
      message: '添加成功',
      type: 'success',
      duration: 1000,
    });
  } else {
    handleAddVocabularyError(addVocabularyResponse);
  }
  const getVocabularyCountResponse = await getVocabularyCount(vocabularySetId);
  checkResponse(getVocabularyCountResponse);
  countTotalRecords.value = getVocabularyCountResponse.data.data;
  const getVocabularyListResponse = await getVocabularyList(vocabularySetId, pageSize.value, currentPage.value);
  checkResponse(getVocabularyListResponse);
  vocabularies = getVocabularyListResponse.data;
  dataLoaded.value = true;
  newVocabulary.vocabulary = '';
  newVocabulary.meaning = '';
  newVocabulary.note = '';
}
</script>

<style scoped>
</style>