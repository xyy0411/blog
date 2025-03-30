import {defineStore} from "pinia";

const blogStore = defineStore('arlog', {
    state: () => ({
      key: "",
      token: "",
      user: "",
      img: "",
      phone: "",
      mail: "",
      gitee: "",
      qq: "",
      wechat: "",
      ip: "",
      otp: "",
      power: ""

    }),
    persist: {
      enabled: true,
      strategies: [
        {
          key: 'arlog',
          storage: localStorage,
        },
      ],
    }
  }
)

export default blogStore
