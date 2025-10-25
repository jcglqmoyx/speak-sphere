<template>
  <ContentBase>
    <el-container v-if="dataLoaded">
      <el-main>
        <el-table :data="dictionaries">
          <el-table-column label="名称" prop="title"/>
          <el-table-column label="前缀" prop="prefix"/>
          <el-table-column label="后缀" prop="suffix"/>
          <el-table-column label="创建时间" prop="created_at" :formatter="formatDateTime"/>
          <el-table-column label="更新时间" prop="updated_at" :formatter="formatDateTime"/>
          <el-table-column label="操作">
            <template #default="scope">
              <el-tooltip content="编辑词典" placement="top">
                <el-button @click="handleEditDictionary(scope.row)">编辑</el-button>
              </el-tooltip>
              <el-tooltip content="删除词典" placement="top">
                <el-button @click="handleDeleteDictionary(scope.row)">删除</el-button>
              </el-tooltip>
            </template>
          </el-table-column>
        </el-table>
      </el-main>
      <br>
      <el-footer>
        <el-tooltip content="添加新词典" placement="top">
          <el-button type="primary" @click="showAddDialog">添加</el-button>
        </el-tooltip>
      </el-footer>
    </el-container>


    <el-dialog v-model="editFormVisible" title="编辑词典">
      <el-form :model="dictionary">
        <el-form-item label="名称">
          <el-input v-model="dictionary.title" autocomplete="off"/>
        </el-form-item>
        <el-form-item label="前缀">
          <el-input v-model="dictionary.prefix" autocomplete="off"/>
        </el-form-item>
        <el-form-item label="后缀">
          <el-input v-model="dictionary.suffix" autocomplete="off"/>
        </el-form-item>
      </el-form>
      <template #footer>
                <span class="dialog-footer">
                  <el-tooltip content="取消编辑" placement="top">
                    <el-button @click="cancelUpdate">取消</el-button>
                  </el-tooltip>
                  <el-tooltip content="确认编辑 (快捷键: Enter)" placement="top">
                    <el-button type="primary" @click="confirmUpdate">确认</el-button>
                  </el-tooltip>
                </span>
      </template>
    </el-dialog>

    <el-dialog v-model="addFormVisible" title="添加词典">
      <el-form :model="newDictionary">
        <el-form-item label="名称">
          <el-input v-model="newDictionary.title" @keyup.enter="confirmAdd" autocomplete="off"/>
        </el-form-item>
        <el-form-item label="前缀">
          <el-input v-model="newDictionary.prefix" @keyup.enter="confirmAdd" autocomplete="off"/>
        </el-form-item>
        <el-form-item label="后缀">
          <el-input v-model="newDictionary.suffix" @keyup.enter="confirmAdd" autocomplete="off"/>
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
import {formatDateTime} from "@/assets/js/util/datetime_util";
import {
  ElButton,
  ElContainer,
  ElDialog,
  ElFooter,
  ElForm,
  ElFormItem,
  ElInput,
  ElMain,
  ElNotification,
  ElTable,
  ElTableColumn,
  ElTooltip,
} from "element-plus";
import {onMounted, reactive, ref} from "vue";
import ContentBase from "@/components/ContentBase.vue";
import {useStore} from "vuex";
import {getDictionaryList} from "@/assets/js/module/dictionary/query";
import {updateDictionary} from "@/assets/js/module/dictionary/update";
import {addDictionary} from "@/assets/js/module/dictionary/add";
import {deleteDictionary} from "@/assets/js/module/dictionary/delete";

let dictionaries = reactive([]);

const store = useStore();
const checkResponse = (response) => {
  if (response == null || response.code === 2) {
    store.dispatch("logout")
    location.reload();
  }
}

const dataLoaded = ref(false);
onMounted(
    async () => {
      const getDictionaryListResponse = await getDictionaryList();
      checkResponse(getDictionaryListResponse);
      dictionaries = getDictionaryListResponse.data;
      dataLoaded.value = true;
    }
);

const dictionary = ref(null);
const editFormVisible = ref(false);


const handleDeleteDictionary = async (row) => {
  ElNotification({
    title: '删除词典',
    message: '词典删除中, 请稍等...',
    type: 'info',
    duration: 1000,
  });
  const deleteDictionaryResponse = await deleteDictionary(row.id);
  checkResponse(deleteDictionaryResponse);
  if (deleteDictionaryResponse.code === 0) {
    dataLoaded.value = false;
    const getDictionaryListResponse = await getDictionaryList();
    checkResponse(getDictionaryListResponse);
    dictionaries = getDictionaryListResponse.data;
    dataLoaded.value = true;
    ElNotification({
      title: 'Success',
      message: '删除成功',
      type: 'success',
      duration: 1000,
    });
  }
}

const handleEditDictionary = (row) => {
  dictionary.value = row;
  editFormVisible.value = true;
}

const cancelUpdate = () => {
  editFormVisible.value = false;
  location.reload();
}
const confirmUpdate = async () => {
  editFormVisible.value = false;
  dataLoaded.value = false;

  const updateDictionaryResponse = await updateDictionary(
      dictionary.value.id,
      dictionary.value.title,
      dictionary.value.prefix,
      dictionary.value.suffix,
  )
  checkResponse(updateDictionaryResponse);
  dataLoaded.value = true;
}

const addFormVisible = ref(false);
const showAddDialog = () => {
  addFormVisible.value = true;
}

const newDictionary = reactive({
  title: '',
  prefix: '',
  suffix: '',
})
const cancelAdd = () => {
  addFormVisible.value = false;
  newDictionary.title = '';
  newDictionary.prefix = '';
  newDictionary.suffix = '';
}
const handleAddDictionaryError = (response) => {
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
  const addDictionaryResponse = await addDictionary(newDictionary.title, newDictionary.prefix, newDictionary.suffix);
  checkResponse(addDictionaryResponse);
  if (addDictionaryResponse.code === 0) {
    ElNotification({
      title: 'Success',
      message: '添加成功',
      type: 'success',
      duration: 1000,
    });
  } else {
    handleAddDictionaryError(addDictionaryResponse);
  }
  const getDictionaryResponse = await getDictionaryList();
  checkResponse(getDictionaryResponse);
  dictionaries = getDictionaryResponse.data;
  dataLoaded.value = true;
  newDictionary.title = '';
  newDictionary.prefix = '';
  newDictionary.suffix = '';
}
</script>

<style scoped>
</style>