import axios from "axios";

const axiosInstance = axios.create({
  baseURL: process.env.REACT_APP_BASE_SERVER_URL,
});

const ApiHandler = {
  client: axiosInstance,
  async searchByIP(address) {
    try {
      const response = await this.client.get(`/ip/${address}`);
      return response.data;
    } catch (error) {
      return { message: [] };
    }
  },
  async searchByDomain(address) {
    console.log(process.env);
    try {
      const response = await this.client.get(`/domain/${address}`);
      return response.data;
    } catch (error) {
      return { message: [] };
    }
  },
  validateIP(str) {
    const regex = /^(([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])\.){3}([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])$/;
    return regex.test(str);
  },
  validateHost(str) {
    const regex = /^(([a-zA-Z0-9]|[a-zA-Z0-9][a-zA-Z0-9-]*[a-zA-Z0-9])\.)*([A-Za-z0-9]|[A-Za-z0-9][A-Za-z0-9-]*[A-Za-z0-9])$/;
    return regex.test(str);
  },
};

export default ApiHandler;
