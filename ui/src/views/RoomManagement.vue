<template>
  <div class="min-h-screen bg-gray-50 dark:bg-gray-900">
    <div class="max-w-7xl mx-auto px-6 py-8">
      <div class="mb-8 flex items-center justify-between">
        <div>
          <router-link to="/" class="text-blue-600 dark:text-blue-400 hover:underline mb-4 inline-block"> ← 返回首页 </router-link>
          <h1 class="text-3xl font-bold text-gray-900 dark:text-white">协作房间</h1>
          <p class="text-gray-600 dark:text-gray-400 mt-2">多 Agent 协作工作空间（演示模式）</p>
        </div>
        <div class="flex items-center gap-3">
          <span class="px-3 py-1 bg-amber-100 dark:bg-amber-900/30 text-amber-700 dark:text-amber-300 text-sm rounded-lg"> 🎭 演示模式 </span>
        </div>
      </div>

      <RoomList :rooms="rooms" :loading="loading" @join="handleJoin" @edit="handleEdit" @delete="handleDelete" />

      <!-- 创建/编辑对话框 -->
      <div v-if="showCreateDialog || showEditDialog" class="fixed inset-0 bg-black/50 flex items-center justify-center z-50" @click.self="closeDialogs">
        <div class="bg-white dark:bg-gray-800 rounded-xl shadow-xl max-w-2xl w-full mx-4 max-h-[90vh] overflow-y-auto">
          <div class="p-6 border-b border-gray-200 dark:border-gray-700">
            <h2 class="text-2xl font-bold text-gray-900 dark:text-white">
              {{ showEditDialog ? "编辑房间" : "创建房间" }}
            </h2>
          </div>

          <div class="p-6 space-y-4">
            <div>
              <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2"> 房间名称 </label>
              <input
                v-model="formData.name"
                type="text"
                placeholder="输入房间名称"
                class="w-full px-4 py-2 border border-gray-300 dark:border-gray-600 rounded-lg bg-white dark:bg-gray-700 text-gray-900 dark:text-white focus:ring-2 focus:ring-blue-500 focus:border-transparent"
              />
            </div>

            <div>
              <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2"> 描述 </label>
              <textarea
                v-model="formData.description"
                rows="3"
                placeholder="输入房间描述"
                class="w-full px-4 py-2 border border-gray-300 dark:border-gray-600 rounded-lg bg-white dark:bg-gray-700 text-gray-900 dark:text-white focus:ring-2 focus:ring-blue-500 focus:border-transparent"
              ></textarea>
            </div>

            <div>
              <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2"> Agents </label>
              <div class="space-y-2">
                <div v-for="(agent, index) in formData.agents" :key="index" class="flex items-center gap-2">
                  <input v-model="formData.agents[index]" type="text" placeholder="Agent 名称" class="flex-1 px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-lg bg-white dark:bg-gray-700 text-gray-900 dark:text-white" />
                  <button @click="removeAgent(index)" class="p-2 text-red-600 hover:bg-red-50 dark:hover:bg-red-900/20 rounded-lg">
                    <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
                    </svg>
                  </button>
                </div>
                <button @click="addAgent" class="w-full px-4 py-2 border-2 border-dashed border-gray-300 dark:border-gray-600 rounded-lg text-gray-600 dark:text-gray-400 hover:border-blue-500 hover:text-blue-600 dark:hover:text-blue-400 transition-colors">
                  + 添加 Agent
                </button>
              </div>
            </div>

            <div>
              <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2"> 最大成员数 </label>
              <input
                v-model.number="formData.maxMembers"
                type="number"
                min="1"
                placeholder="最大成员数"
                class="w-full px-4 py-2 border border-gray-300 dark:border-gray-600 rounded-lg bg-white dark:bg-gray-700 text-gray-900 dark:text-white focus:ring-2 focus:ring-blue-500 focus:border-transparent"
              />
            </div>
          </div>

          <div class="p-6 border-t border-gray-200 dark:border-gray-700 flex justify-end gap-3">
            <button @click="closeDialogs" class="px-4 py-2 border border-gray-300 dark:border-gray-600 text-gray-700 dark:text-gray-300 rounded-lg hover:bg-gray-50 dark:hover:bg-gray-700 transition-colors">取消</button>
            <button @click="showEditDialog ? updateRoom() : createRoom()" class="px-4 py-2 bg-blue-600 hover:bg-blue-700 text-white rounded-lg transition-colors">
              {{ showEditDialog ? "保存" : "创建" }}
            </button>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, reactive } from "vue";
import RoomList from "../components/Room/RoomList.vue";
import { useNomadClient } from "../composables/useNomadClient";
import { DEMO_MODE, demoRooms } from "../config/demoData";

const { client } = useNomadClient();
const rooms = ref<any[]>([]);

console.log("🎭 演示模式:", DEMO_MODE ? "启用（使用本地数据）" : "禁用（连接后端API）");

const loading = ref(false);
const showCreateDialog = ref(false);
const showEditDialog = ref(false);
const editingRoomId = ref<string | null>(null);

const formData = reactive({
  name: "",
  description: "",
  agents: [""],
  maxMembers: 10,
});

// Create - 演示模式禁用创建
const createRoom = async () => {
  alert("演示模式下不支持创建新房间\n\n您可以查看和加入现有的演示房间");
  closeDialogs();
};

// Read - 加载房间列表
const loadRooms = async () => {
  try {
    loading.value = true;

    if (DEMO_MODE) {
      // 演示模式：使用 UI 本地数据
      rooms.value = JSON.parse(JSON.stringify(demoRooms));
    } else {
      // 生产模式：从后端 API 获取
      const response = await client.rooms.list();
      if (response.success && response.data) {
        rooms.value = response.data.map((r: any) => ({
          ...r,
          description: r.metadata?.description || "",
          agents: r.metadata?.agents || [],
          members: r.metadata?.members || 0,
          maxMembers: r.metadata?.maxMembers || 10,
          status: "active",
        }));
      }
    }
  } catch (error: any) {
    console.error("加载房间失败:", error);
    // 失败时使用演示数据作为后备
    rooms.value = JSON.parse(JSON.stringify(demoRooms));
  } finally {
    loading.value = false;
  }
};

// Update - 更新房间（后端暂不支持，使用本地更新）
const updateRoom = () => {
  if (!formData.name.trim()) {
    alert("请输入房间名称");
    return;
  }

  const index = rooms.value.findIndex((r) => r.id === editingRoomId.value);
  if (index !== -1) {
    rooms.value[index] = {
      ...rooms.value[index],
      name: formData.name,
      description: formData.description,
      agents: formData.agents.filter((a) => a.trim()),
      maxMembers: formData.maxMembers,
    };
    closeDialogs();
    alert(`房间 "${formData.name}" 更新成功！`);
  }
};

// Delete - 演示模式禁用删除
const handleDelete = async (room: any) => {
  alert("演示模式下不支持删除房间\n\n这是一个只读演示环境");
};

// Join - 加入房间（演示模式：模拟加入）
const handleJoin = async (room: any) => {
  const index = rooms.value.findIndex((r) => r.id === room.id);
  if (index !== -1) {
    if (rooms.value[index].members < rooms.value[index].maxMembers) {
      rooms.value[index].members++;
      alert(`成功加入房间: ${room.name}\n\n当前成员: ${rooms.value[index].members} 人\nAgents: ${room.agents.join(", ")}\n\n演示模式：这是一个模拟操作`);
    } else {
      alert("房间已满，无法加入");
    }
  }
};

// Edit - 演示模式禁用编辑
const handleEdit = (room: any) => {
  alert("演示模式下不支持编辑房间\n\n您可以查看房间详情和加入演示");
};

// Agent 管理
const addAgent = () => {
  formData.agents.push("");
};

const removeAgent = (index: number) => {
  if (formData.agents.length > 1) {
    formData.agents.splice(index, 1);
  } else {
    alert("至少需要保留一个 Agent");
  }
};

const closeDialogs = () => {
  showCreateDialog.value = false;
  showEditDialog.value = false;
  editingRoomId.value = null;
  formData.name = "";
  formData.description = "";
  formData.agents = [""];
  formData.maxMembers = 10;
};

onMounted(() => {
  loadRooms();
});
</script>
