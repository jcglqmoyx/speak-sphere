<template>
  <ContentBase>
    <el-container v-if="dataLoaded">
      <el-header>
        <el-button type="primary" @click="showAddDialog">添加AI提示词</el-button>
      </el-header>
      <el-main>
        <el-table :data="prompts" empty-text="暂无AI提示词配置">
          <el-table-column label="提示词名称" prop="name"/>
          <el-table-column label="描述" prop="description"/>
          <el-table-column label="是否默认">
            <template #default="scope">
              <el-tag v-if="scope.row.is_default" type="success">默认</el-tag>
              <el-button v-else type="text" @click="setDefaultPrompt(scope.row)">设为默认</el-button>
            </template>
          </el-table-column>
          <el-table-column label="创建时间" prop="created_at" :formatter="formatDateTime"/>
          <el-table-column label="操作">
            <template #default="scope">
              <el-button v-if="canEdit(scope.row)" @click="handleEditPrompt(scope.row)">编辑</el-button>
              <el-button v-if="canEdit(scope.row)" @click="handleDeletePrompt(scope.row)" type="danger">删除</el-button>
            </template>
          </el-table-column>
        </el-table>
      </el-main>
    </el-container>

    <!-- 添加/编辑提示词对话框 -->
    <el-dialog v-model="formVisible" :title="formTitle" width="800px">
      <el-form :model="promptForm" :rules="formRules" ref="promptFormRef">
        <el-form-item label="提示词名称" prop="name">
          <el-input v-model="promptForm.name" placeholder="请输入提示词名称"/>
        </el-form-item>
        <el-form-item label="提示词描述">
          <el-input 
            v-model="promptForm.description" 
            placeholder="请输入提示词描述（可选）"
          />
        </el-form-item>
        <el-form-item label="提示词内容" prop="content">
          <el-input 
            v-model="promptForm.content" 
            type="textarea"
            placeholder="请输入提示词内容"
            :rows="8"
            resize="vertical"
          />
        </el-form-item>
        <el-form-item>
          <el-checkbox v-model="promptForm.is_default">设为默认提示词</el-checkbox>
        </el-form-item>
      </el-form>
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="cancelForm">取消</el-button>
          <el-button type="primary" @click="confirmForm">确认</el-button>
        </span>
      </template>
    </el-dialog>
  </ContentBase>
</template>

<script setup>
import {computed, nextTick, onMounted, reactive, ref} from "vue";
import {addAIPrompt, deleteAIPrompt, queryAIPrompts, updateAIPrompt} from "@/assets/js/module/aiprompt";
import {
  ElButton,
  ElCheckbox,
  ElContainer,
  ElDialog,
  ElForm,
  ElFormItem,
  ElHeader,
  ElInput,
  ElMain,
  ElMessage,
  ElMessageBox,
  ElNotification,
  ElTable,
  ElTableColumn,
  ElTag,
} from "element-plus";
import "element-plus/dist/index.css";
import {useStore} from "vuex";
import ContentBase from "@/components/ContentBase.vue";
import {formatDateTime} from "@/assets/js/util/datetime_util";

const store = useStore();

// 检查响应状态，处理未授权
const checkResponse = (response) => {
  if (response == null || response.code === 2) {
    store.dispatch("logout")
    location.reload();
  }
}

// 响应式数据
const dataLoaded = ref(false);
const prompts = reactive([]);
const formVisible = ref(false);
const isEditing = ref(false);
const promptFormRef = ref(null);

// 表单数据
const promptForm = reactive({
  id: null,
  name: '',
  description: '',
  content: '',
  is_default: false
});

// 表单验证规则
const formRules = {
  name: [{ required: true, message: '请输入提示词名称', trigger: 'blur' }],
  content: [{ required: true, message: '请输入提示词内容', trigger: 'blur' }]
};

// 计算属性
const formTitle = computed(() => isEditing.value ? '编辑AI提示词' : '添加AI提示词');

// 检查是否可以编辑（只能编辑用户自定义的提示词）
const canEdit = (prompt) => {
  return prompt.user_id !== 0; // 系统默认提示词不可编辑删除
}

// 生命周期
onMounted(async () => {
  await loadPrompts();
});

// 加载提示词列表
const loadPrompts = async () => {
  dataLoaded.value = false;
  try {
    const response = await queryAIPrompts();
    checkResponse(response);
    if (response.code === 0) {
      prompts.splice(0, prompts.length, ...response.data);
    }
  } catch (error) {
    ElNotification({
      title: '错误',
      message: '加载提示词列表失败',
      type: 'error',
    });
  } finally {
    dataLoaded.value = true;
  }
};

// 显示添加对话框
const showAddDialog = () => {
  isEditing.value = false;
  resetForm();
  formVisible.value = true;
};

// 编辑提示词
const handleEditPrompt = (prompt) => {
  if (!canEdit(prompt)) {
    ElMessage.warning('系统默认提示词不可编辑');
    return;
  }
  isEditing.value = true;
  Object.assign(promptForm, prompt);
  formVisible.value = true;
};

// 删除提示词
const handleDeletePrompt = async (prompt) => {
  if (!canEdit(prompt)) {
    ElMessage.warning('系统默认提示词不可删除');
    return;
  }
  
  try {
    await ElMessageBox.confirm(
      `确认删除提示词 "${prompt.name}" 吗？`,
      '删除确认',
      {
        confirmButtonText: '确认',
        cancelButtonText: '取消',
        type: 'warning',
      }
    );
    
    const response = await deleteAIPrompt(prompt.id);
    checkResponse(response);
    
    if (response.code === 0) {
      ElNotification({
        title: '成功',
        message: '提示词删除成功',
        type: 'success',
      });
      await loadPrompts();
    }
  } catch (error) {
    if (error !== 'cancel') {
      ElNotification({
        title: '错误',
        message: '删除提示词失败',
        type: 'error',
      });
    }
  }
};

// 设置默认提示词
const setDefaultPrompt = async (prompt) => {
  try {
    const response = await updateAIPrompt(prompt.id, { is_default: true });
    checkResponse(response);
    
    if (response.code === 0) {
      ElNotification({
        title: '成功',
        message: '默认提示词设置成功',
        type: 'success',
      });
      await loadPrompts();
    }
  } catch (error) {
    ElNotification({
      title: '错误',
      message: '设置默认提示词失败',
      type: 'error',
    });
  }
};

// 重置表单
const resetForm = () => {
  Object.assign(promptForm, {
    id: null,
    name: '',
    description: '',
    content: '',
    is_default: false
  });
  nextTick(() => {
    promptFormRef.value?.clearValidate();
  });
};

// 取消表单
const cancelForm = () => {
  formVisible.value = false;
  resetForm();
};

// 确认表单提交
const confirmForm = async () => {
  if (!promptFormRef.value) return;
  
  try {
    await promptFormRef.value.validate();
    
    ElNotification({
      title: '处理中',
      message: `${isEditing.value ? '更新' : '添加'}提示词中...`,
      type: 'info',
    });
    
    let response;
    if (isEditing.value) {
      // 更新提示词
      response = await updateAIPrompt(promptForm.id, promptForm);
    } else {
      // 添加提示词
      response = await addAIPrompt(promptForm);
    }
    
    checkResponse(response);
    
    if (response.code === 0) {
      ElNotification({
        title: '成功',
        message: `${isEditing.value ? '更新' : '添加'}提示词成功`,
        type: 'success',
      });
      formVisible.value = false;
      resetForm();
      await loadPrompts();
    }
  } catch (error) {
    if (error.errors) {
      ElMessage.warning('请完善表单信息');
    } else {
      ElNotification({
        title: '错误',
        message: `${isEditing.value ? '更新' : '添加'}提示词失败`,
        type: 'error',
      });
    }
  }
};
</script>

<style scoped>
.dialog-footer button:first-child {
  margin-right: 10px;
}
</style>
