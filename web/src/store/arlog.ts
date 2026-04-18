import { defineStore } from 'pinia';

interface BlogState {
  key: string;
  token: string;
  user: string;
  img: string;
  phone: string;
  mail: string;
  gitee: string;
  qq: string;
  wechat: string;
  ip: string;
  otp: string;
  power: string;
}

const defaultState = (): BlogState => ({
  key: '',
  token: '',
  user: '',
  img: '',
  phone: '',
  mail: '',
  gitee: '',
  qq: '',
  wechat: '',
  ip: '',
  otp: '',
  power: '',
});

const loadState = (): BlogState => {
  if (typeof window === 'undefined') {
    return defaultState();
  }

  const saved = window.localStorage.getItem('arlog');
  if (!saved) {
    return defaultState();
  }

  try {
    return {
      ...defaultState(),
      ...JSON.parse(saved),
    };
  } catch {
    return defaultState();
  }
};

const blogStore = defineStore('arlog', {
  state: loadState,
  actions: {
    persist() {
      if (typeof window === 'undefined') {
        return;
      }

      window.localStorage.setItem('arlog', JSON.stringify(this.$state));
    },
    patchAndPersist(payload: Partial<BlogState>) {
      this.$patch(payload);
      this.persist();
    },
    clearAndPersist() {
      this.$reset();
      this.persist();
    },
  },
});

export default blogStore;
