import axios from "axios";
import base from "@/api/api.ts";


const st = blogStore()

/*创建请求对象*/
const req = axios.create(
  {
    baseURL: base.url,
    /*线上修改成请求当前域名*/

    timeout: 5000,
  }
)
/*注册请求拦截*/

/*添加请求token*/
/*所有post请求添加成json格式*/
req.interceptors.request.use((config) => {

    config.headers['token'] = st.token

    if (config.method === 'POST') {
      config.headers.append('Content-Type', 'application/json');
    }
    return config
  },
  error => Promise.reject(error)
);


// 响应拦截器：统一处理错误响应
// 添加响应拦截器
req.interceptors.response.use(res => {

    // 对响应数据做点什么
    const data = res.data
    // 访问成功 清理 之前的请求
    $message.destroy("loadingMessage")
    if (data.code !== 200) {
      // 统一处理非200的状态码
      window.res.error(data.msg)
      return
    }
    return res
  }, error => Promise.reject(error)
)


export default req
