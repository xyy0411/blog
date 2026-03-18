import axios from 'axios';
import base from '@/api/api.ts';
import blogStore from '@/store/arlog.ts';

const req = axios.create({
  baseURL: base.url,
  timeout: 5000,
});

req.interceptors.request.use(
  (config) => {
    const store = blogStore();

    if (store.token) {
      config.headers.token = store.token;
    }

    if (config.method?.toUpperCase() === 'POST') {
      config.headers['Content-Type'] = 'application/json';
    }

    return config;
  },
  (error) => Promise.reject(error),
);

req.interceptors.response.use(
  (response) => {
    const data = response.data;

    if (data?.code && data.code !== 200) {
      return Promise.reject(new Error(data.msg || '请求失败'));
    }

    return response;
  },
  (error) => Promise.reject(error),
);

export default req;
