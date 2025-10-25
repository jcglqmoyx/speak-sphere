<template>
  <div id="app">
    <div class="container">
      <h1>录音机</h1>
      <div class="controls">
        <button class="btn" @click="startRecording">开始录音</button>
        <button class="btn" @click="stopRecording">停止录音</button>
        <input class="input" type="number" v-model="duration" placeholder="自动停止时间（秒）"/>
      </div>
      <input class="input" v-if="countdown > 0" v-model="countdown" placeholder="剩余时间" readonly/>
      <audio v-if="audioUrl" :src="audioUrl" controls></audio>
    </div>
  </div>
</template>

<script>
export default {
  data() {
    return {
      mediaRecorder: null,
      audioChunks: [],
      audioUrl: null,
      countdown: 0, // 倒计时的初始值
      countdownInterval: null, // 倒计时的Interval引用
      duration: 45, // 用户输入的录音持续时间，默认为45秒
    };
  },
  methods: {
    async startRecording() {
      try {
        const stream = await navigator.mediaDevices.getUserMedia({audio: true});
        this.mediaRecorder = new MediaRecorder(stream);
        this.audioChunks = [];

        this.mediaRecorder.ondataavailable = event => {
          this.audioChunks.push(event.data);
        };

        this.mediaRecorder.start();
        this.startCountdown(this.duration); // 使用用户输入的时间作为倒计时
      } catch (err) {
        console.error('Error accessing audio', err);
      }
    },
    stopRecording() {
      this.mediaRecorder.stop();
      this.mediaRecorder.onstop = () => {
        const audioBlob = new Blob(this.audioChunks);
        this.audioUrl = URL.createObjectURL(audioBlob);
      };
      this.stopCountdown();
    },
    startCountdown(duration) {
      this.countdown = duration;
      this.countdownInterval = setInterval(() => {
        this.countdown--;
        if (this.countdown <= 0) {
          this.stopRecording();
        }
      }, 1000);
    },
    stopCountdown() {
      clearInterval(this.countdownInterval);
      this.countdown = 0; // 重置倒计时
    },
  },
};
</script>

<style scoped>
.container {
  display: flex;
  flex-direction: column;
  align-items: center;
  padding: 20px;
}

.controls {
  display: flex;
  gap: 10px;
  margin-bottom: 10px;
}

.btn {
  background-color: #4CAF50;
  color: white;
  padding: 10px 15px;
  border: none;
  cursor: pointer;
  border-radius: 5px;
}

.btn:hover {
  background-color: #45a049;
}

.input {
  padding: 10px;
  border-radius: 5px;
  border: 1px solid #ccc;
}

.input:focus {
  outline: none;
  border-color: #4CAF50;
}

audio {
  margin-top: 10px;
}
</style>