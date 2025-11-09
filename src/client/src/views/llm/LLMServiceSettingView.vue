<template>
  <ContentBase>
    <el-container v-if="dataLoaded">
      <el-header>
        <el-button type="primary" @click="showAddDialog">添加LLM服务</el-button>
      </el-header>
      <el-main>
        <el-table :data="services" empty-text="暂无LLM服务配置">
          <el-table-column label="服务名称" prop="name"/>
          <el-table-column label="提供商" prop="provider"/>
          <el-table-column label="API端点" prop="endpoint"/>
          <el-table-column label="模型" prop="model"/>
          <el-table-column label="是否默认">
            <template #default="scope">
              <el-tag v-if="scope.row.is_default" type="success">默认</el-tag>
              <el-button v-else type="text" @click="setDefaultService(scope.row)">设为默认</el-button>
            </template>
          </el-table-column>
          <el-table-column label="描述" prop="description"/>
          <el-table-column label="创建时间" prop="created_at" :formatter="formatDateTime"/>
          <el-table-column label="操作">
            <template #default="scope">
              <el-button @click="handleEditService(scope.row)">编辑</el-button>
              <el-button @click="handleDeleteService(scope.row)" type="danger">删除</el-button>
            </template>
          </el-table-column>
        </el-table>
      </el-main>
    </el-container>

    <!-- 添加/编辑服务对话框 -->
    <el-dialog v-model="formVisible" :title="formTitle">
      <el-form :model="serviceForm" :rules="formRules" ref="serviceFormRef">
        <el-form-item label="服务名称" prop="name">
          <el-input v-model="serviceForm.name" placeholder="请输入服务名称"/>
        </el-form-item>
        <el-form-item label="服务提供商">
          <el-input v-model="serviceForm.provider" placeholder="请输入服务提供商"/>
        </el-form-item>
        <el-form-item label="API端点" prop="endpoint">
          <el-input v-model="serviceForm.endpoint" placeholder="请输入API端点URL"/>
        </el-form-item>
        <el-form-item label="模型名称" prop="model">
          <el-input v-model="serviceForm.model" placeholder="请输入模型名称"/>
        </el-form-item>
        <el-form-item label="API Key" prop="api_key">
          <el-input
              v-model="serviceForm.api_key"
              type="password"
              placeholder="请输入API Key"
              show-password
          />
        </el-form-item>
        <el-form-item label="服务描述">
          <el-input
              v-model="serviceForm.description"
              type="textarea"
              placeholder="请输入服务描述"
              :rows="3"
          />
        </el-form-item>
        <el-form-item>
          <el-checkbox v-model="serviceForm.is_default">设为默认服务</el-checkbox>
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
import {addLLMService, deleteLLMService, queryLLMServices, updateLLMService} from "@/assets/js/module/llm";
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
const services = reactive([]);
const formVisible = ref(false);
const isEditing = ref(false);
const serviceFormRef = ref(null);

// 表单数据
const serviceForm = reactive({
  id: null,
  name: '',
  provider: '',
  endpoint: '',
  model: '',
  api_key: '',
  description: '',
  is_default: false
});

// 表单验证规则
const formRules = {
  name: [{required: true, message: '请输入服务名称', trigger: 'blur'}],
  endpoint: [{required: true, message: '请输入API端点', trigger: 'blur'}],
  model: [{required: true, message: '请输入模型名称', trigger: 'blur'}],
  api_key: [{required: true, message: '请输入API Key', trigger: 'blur'}]
};

// 计算属性
const formTitle = computed(() => isEditing.value ? '编辑LLM服务' : '添加LLM服务');

// 生命周期
onMounted(async () => {
  await loadServices();
});

// 加载服务列表
const loadServices = async () => {
  dataLoaded.value = false;
  try {
    const response = await queryLLMServices();
    checkResponse(response);
    if (response.code === 0) {
      services.splice(0, services.length, ...response.data);
    }
  } catch (error) {
    ElNotification({
      title: '错误',
      message: '加载服务列表失败',
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

// 编辑服务
const handleEditService = (service) => {
  isEditing.value = true;
  Object.assign(serviceForm, service);
  formVisible.value = true;
};

// 删除服务
const handleDeleteService = async (service) => {
  try {
    await ElMessageBox.confirm(
        `确认删除服务 "${service.name}" 吗？`,
        '删除确认',
        {
          confirmButtonText: '确认',
          cancelButtonText: '取消',
          type: 'warning',
        }
    );

    const response = await deleteLLMService(service.id);
    checkResponse(response);

    if (response.code === 0) {
      ElNotification({
        title: '成功',
        message: '服务删除成功',
        type: 'success',
      });
      await loadServices();
    }
  } catch (error) {
    if (error !== 'cancel') {
      ElNotification({
        title: '错误',
        message: '删除服务失败',
        type: 'error',
      });
    }
  }
};

// 设置默认服务
const setDefaultService = async (service) => {
  try {
    const response = await updateLLMService(service.id, {is_default: true});
    checkResponse(response);

    if (response.code === 0) {
      ElNotification({
        title: '成功',
        message: '默认服务设置成功',
        type: 'success',
      });
      await loadServices();
    }
  } catch (error) {
    ElNotification({
      title: '错误',
      message: '设置默认服务失败',
      type: 'error',
    });
  }
};

// 重置表单
const resetForm = () => {
  Object.assign(serviceForm, {
    id: null,
    name: '',
    provider: '',
    endpoint: '',
    model: '',
    api_key: '',
    description: '',
    is_default: false
  });
  nextTick(() => {
    serviceFormRef.value?.clearValidate();
  });
};

// 取消表单
const cancelForm = () => {
  formVisible.value = false;
  resetForm();
};

// 确认表单提交
const confirmForm = async () => {
  if (!serviceFormRef.value) return;

  try {
    await serviceFormRef.value.validate();

    ElNotification({
      title: '处理中',
      message: `${isEditing.value ? '更新' : '添加'}服务中...`,
      type: 'info',
    });

    let response;
    if (isEditing.value) {
      // 更新服务
      response = await updateLLMService(serviceForm.id, serviceForm);
    } else {
      // 添加服务
      response = await addLLMService(serviceForm);
    }

    checkResponse(response);

    if (response.code === 0) {
      ElNotification({
        title: '成功',
        message: `${isEditing.value ? '更新' : '添加'}服务成功`,
        type: 'success',
      });
      formVisible.value = false;
      resetForm();
      await loadServices();
    }
  } catch (error) {
    if (error.errors) {
      ElMessage.warning('请完善表单信息');
    } else {
      ElNotification({
        title: '错误',
        message: `${isEditing.value ? '更新' : '添加'}服务失败`,
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