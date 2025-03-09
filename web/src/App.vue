<template>
  <body>
  <section id = "login-form">
    <h1>WELCOME</h1>

    <div class="input-wrap">
      <input  v-model="emails" type="text" placeholder="Your Email"
             spellcheck="false" required>

      <i class="fa-solid fa-user"></i>
    </div>

    <div class="input-wrap">
      <input v-model="password" type="password" placeholder="Your Password"
             spellcheck="false" required>

      <i class="fa-solid fa-lock"></i>
    </div>

    <div class="rem">
      <p>
        <input type="checkbox">
        记住我
      </p>
      <a @click="findPassword">Forgot password?</a>
    </div>

    <button @click="login">Login</button>

    <p class="reg" @click="registerUser">
      Don't have an account?
      <a>register</a>
    </p>
  </section>
  </body>
</template>
<script setup lang="ts">
  import {ref} from "vue";
  import axios from "axios";

  let emails = ref(``);
  let password = ref(``);

  function findPassword() {
    alert(`找回`)
  }

  function registerUser() {

  }


  async function login() {
    const url = `http://127.0.0.1:3000/api/auth/login`;

    if (emails.value === '' || password.value === '') {
      alert("账号或密码未输入");
      return;
    }

    try {
      const response = await axios.post(url, {
        name: emails.value,
        password: password.value
      });
      console.log(response);
    } catch (error) {
      console.error('Error fetching data:', error);
    }
  }


</script>

<style scoped>

  * {
    padding: 0;
    margin: 0;
    outline: none;
    border: none;
  }

  body {
    background-image: url("public/login.png");
    background-size: cover;
    background-position: center;
    height: 100vh;
    color: white;

    display: flex;
    justify-content: center;
    align-items: center;
  }

  #login-form {
    display: flex;
    justify-content: center;
    align-items: center;
    flex-direction: column;
    background-color: rgb(255,255,255,0.1);
    width: 35%;
    padding: 50px;
    border-radius: 10px;
    border: 1px solid rgba(255,255,255,0.329);
    backdrop-filter: blur(10px);
    box-shadow: 0 10px 25px rgba(0,0,0,0.13);
  }

  #login-form {
    padding: 50px 10px;
  }

  h1 {
    font-size: 2.5rem;
    color: white;
    margin: 0 0 50px;
  }

  .input-wrap {
    border: 1px solid white;
    display: flex;
    align-items: center;
    justify-content: space-between;
    margin-bottom: 40px;
    height: 70%;
    width: 90%;
    border-radius: 50px;
  }

  input {
    background-color: transparent;
    font-size: 2rem;
    color: white;
    padding: 2px 25px;
  }

  input::placeholder {
    background-color: transparent;
    font-size: 1.3rem;
    color: white;
  }

  .input-wrap i {
    background-color: transparent;
    color: white;
    padding-right: 25px ;
  }

  .rem {
    display: flex;
    justify-content: space-between;
    align-items: center;
    width: 90%;
    margin-bottom: 40px;
  }

  .rem p,.rem a {
    font-size: 1.2rem;
    color: white;

    cursor: pointer;
  }

  button {
    font-size: 1.2rem;
    height: 60px;
    width: 90%;
    border-radius: 50px ;
    font-weight: 600;
    letter-spacing: 1px;
    margin-bottom: 25px;
  }

  button:hover {
    background-color: rgb(211,211,211);
  }

  .reg {
    font-size: 1.2rem;
    color: white;
    margin-top: 40px;
    cursor: pointer;
    transition: 0.3s;
  }

  .reg a {
    font-weight: 500;
    cursor: pointer;
  }
</style>
