<template>
  <div id="login-form">
    <h1>注册</h1>
    <div class="input-wrap">
      <input v-model="username" placeholder="用户名" />
    </div>
    <div class="input-wrap">
      <input v-model="password" type="password" placeholder="密码" />
    </div>
    <div class="input-wrap">
      <input v-model="email" placeholder="邮箱" />
    </div>
    <button @click="register">注册</button>
    <p class="reg">
      已有账号？<a @click="goToLogin">登录</a>
    </p>
  </div>
</template>

<script lang="ts" setup>
import { ref } from 'vue';
import req from '@/utils/requset.ts';
import { useRouter } from 'vue-router';

// 初始化路由实例
const router = useRouter();
// 定义用户名、密码和邮箱的响应式变量，并添加类型注解
const username = ref<string>('');
const password = ref<string>('');
const email = ref<string>('');

/**
 * 处理用户注册逻辑
 */
const register = async () => {
  try {
    // 发起注册请求
    const response = await req.post<{ message: string }>('/api/auth/register', {
      username: username.value,
      password: password.value,
      email: email.value,
    });
    console.log('注册成功', response.data);
    // 注册成功后跳转到登录页面
    router.push('/login');
  } catch (error) {
    console.error('注册失败', error);
    // 给用户友好的错误提示
    alert('注册失败，请检查输入信息');
  }
};

/**
 * 跳转到登录页面
 */
const goToLogin = () => {
  router.push('/login');
};
</script>

<style scoped>
/* 整体表单容器样式 */
#login-form {
  width: 300px;
  margin: 50px auto;
  padding: 20px;
  border: 1px solid #ccc;
  border-radius: 5px;
  box-shadow: 0 0 10px rgba(0, 0, 0, 0.1);
  background-color: #fff;
}

/* 标题样式 */
#login-form h1 {
  text-align: center;
  margin-bottom: 20px;
  color: #333;
}

/* 输入框包裹层样式 */
.input-wrap {
  margin-bottom: 15px;
}

/* 输入框样式 */
.input-wrap input {
  width: 100%;
  padding: 10px;
  border: 1px solid #ccc;
  border-radius: 3px;
  box-sizing: border-box;
}

/* 按钮样式 */
#login-form button {
  width: 100%;
  padding: 10px;
  background-color: #007bff;
  color: #fff;
  border: none;
  border-radius: 3px;
  cursor: pointer;
}

/* 按钮悬停效果 */
#login-form button:hover {
  background-color: #0056b3;
}

/* 注册提示文字样式 */
.reg {
  text-align: center;
  margin-top: 15px;
}

/* 链接样式 */
.reg a {
  color: #007bff;
  text-decoration: none;
}

/* 链接悬停效果 */
.reg a:hover {
  text-decoration: underline;
}
</style>

