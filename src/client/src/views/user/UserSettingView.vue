<template>
  <ContentBase>
    <el-card v-if="dataLoaded">
      <el-row>
        <el-col :span="9">
          <el-avatar :size="120" :src="user.avatar"></el-avatar>
        </el-col>
        <el-col :span="16">
          <el-row>
            <el-col :span="12">
              <h2>{{ user.username }}</h2>
            </el-col>
          </el-row>
          <el-row>
            <el-col :span="24">
              <p>Email: <b>{{ user.email }}</b></p>
              <p>当前所学语书: <b>{{ vocabularySetToLearn.title }}</b></p>
              <p>每日新单词数: <b>{{ user.dailyCount }}</b></p>
              <p>复习频率公式: <b>{{ user.reviewFrequencyFormula }}</b></p>
              <p>点击几次“认识”算学会: <b>{{ user.timesCountedAsKnown }}</b></p>
              <p>
                主题:
                <el-switch
                    v-model="isDark"
                    class="mt-2"
                    style="margin-left: 24px"
                    inline-prompt
                    :active-icon="Sunny"
                    :inactive-icon="Moon"
                />
              </p>
              <el-button @click="handleEdit()">修改</el-button>
            </el-col>
          </el-row>
        </el-col>
      </el-row>
    </el-card>

    <el-dialog v-model="editFormVisible" title="修改设置">
      <el-form :model="user">
        <el-form-item label="用户名">
          <el-input v-model="user.username"/>
        </el-form-item>
        <el-form-item label="邮箱">
          <el-input v-model="user.email"/>
        </el-form-item>
        <el-form-item label="头像">
          <el-input v-model="user.avatar"/>
        </el-form-item>
        <el-form-item label="当前所学词书">
          <el-dropdown split-button type="primary">
            {{ vocabularySetToLearn.title || "选择词书" }}
            <template #dropdown>
              <el-dropdown-menu>
                <el-dropdown-item v-for="vocabularySet in vocabularySets" :key="vocabularySet.id"
                                  @click="handleChooseVocabularySetToLearn(vocabularySet.id)">{{
                    vocabularySet.title
                  }}
                </el-dropdown-item>
              </el-dropdown-menu>
            </template>
          </el-dropdown>

        </el-form-item>
        <el-form-item label="每日新单词数">
          <el-input-number v-model="user.dailyCount" :min="1" :max="1000"/>
        </el-form-item>
        <el-form-item label="复习频率公式">
          <el-input v-model="user.reviewFrequencyFormula"/>
        </el-form-item>
        <el-form-item label="点击几次“认识”算学会">
          <el-input-number v-model="user.timesCountedAsKnown" :min="1" :max="1000"/>
        </el-form-item>
      </el-form>
      <template #footer>
                <span class="dialog-footer">
                  <el-button @click="cancelUpdate">取消</el-button>
                  <el-button type="primary" @click="confirmUpdate">确认</el-button>
                </span>
      </template>
    </el-dialog>
  </ContentBase>
</template>

<script setup>
import "element-plus/dist/index.css"
import {onMounted, reactive, ref} from "vue";
import {
  ElAvatar,
  ElButton,
  ElCard,
  ElCol,
  ElDialog,
  ElDropdown,
  ElDropdownItem,
  ElDropdownMenu,
  ElForm,
  ElFormItem,
  ElInput,
  ElInputNumber,
  ElMessage,
  ElRow,
  ElSwitch,
} from "element-plus";
import ContentBase from "@/components/ContentBase.vue";
import {getUserProfile} from "@/assets/js/module/user/query";
import {useStore} from "vuex";
import {updateUser} from "@/assets/js/module/user/update";
import {getVocabularySetList} from "@/assets/js/module/vocabulary_set/query";
import {getThemeInstance} from "@/store/theme.js";
import {Moon, Sunny} from "@element-plus/icons-vue";

let vocabularySets = reactive([]);
let vocabularySetToLearn = reactive({
  id: 0,
  title: "",
});

const handleChooseVocabularySetToLearn = (id) => {
  for (let vocabularySet of vocabularySets) {
    if (vocabularySet.id === id) {
      vocabularySetToLearn.id = vocabularySet.id;
      vocabularySetToLearn.title = vocabularySet.title;
      user.currentVocabularySetID = vocabularySet.id;
      break;
    }
  }
}
let user = reactive(
    {
      username: "",
      email: "",
      avatar: "",
      currentVocabularySetID: 0,
      dailyCount: 0,
      timesCountedAsKnown: 0,
      reviewFrequencyFormula: "",
    }
)
const store = useStore();
const dataLoaded = ref(false);
const checkResponse = (response) => {
  if (response == null || response.code === 2) {
    store.dispatch("logout")
    location.reload();
  }
}
onMounted(
    async () => {
      dataLoaded.value = false;
      const getUserProfileResponse = await getUserProfile();
      checkResponse(getUserProfileResponse);
      user.username = getUserProfileResponse.username;
      user.email = getUserProfileResponse.email;
      user.avatar = getUserProfileResponse.avatar;
      user.currentVocabularySetID = getUserProfileResponse.current_vocabularySet_id;
      user.dailyCount = getUserProfileResponse.daily_count;
      user.timesCountedAsKnown = getUserProfileResponse.times_counted_as_known;
      user.reviewFrequencyFormula = getUserProfileResponse.review_frequency_formula;

      const getVocabularySetListResponse = await getVocabularySetList(100000000, 1);
      checkResponse(getVocabularySetListResponse);
      vocabularySets = getVocabularySetListResponse.data;
      for (let vocabularySet of vocabularySets) {
        if (vocabularySet.id === user.currentVocabularySetID) {
          vocabularySetToLearn.id = vocabularySet.id;
          vocabularySetToLearn.title = vocabularySet.title;
          break;
        }
      }

      dataLoaded.value = true;
    }
)

const editFormVisible = ref(false);
const handleEdit = () => {
  editFormVisible.value = true;
}
const cancelUpdate = () => {
  editFormVisible.value = false;
  location.reload();
}
const confirmUpdate = async () => {
  editFormVisible.value = false;
  dataLoaded.value = false;

  const updateUserResponse = await updateUser(
      user.username,
      user.email,
      user.avatar,
      +user.currentVocabularySetID,
      user.dailyCount,
      user.timesCountedAsKnown,
      user.reviewFrequencyFormula,
  )

  if (updateUserResponse && updateUserResponse.code === 0) {
    checkResponse(updateUserResponse);
    dataLoaded.value = true;
    // 显示成功消息
    ElMessage.success('用户信息更新成功');
    // 延迟一点时间再刷新，让用户看到成功消息
    setTimeout(() => {
      location.reload();
    }, 1000);
  } else {
    // 更新失败，重新显示编辑表单并显示错误信息
    editFormVisible.value = true;
    dataLoaded.value = true;
    const errorMsg = updateUserResponse ? updateUserResponse.message : '更新失败，请重试';
    ElMessage.error(errorMsg);
  }
}

const {isDark} = getThemeInstance();
</script>
<style scoped>
</style>