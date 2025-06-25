import axios from 'axios';

const API_BASE = 'http://localhost:8080';

export const listScripts = async () => {
  const res = await axios.get(`${API_BASE}/admin/scripts`);
  return res.data;
};

export const saveScript = async (data: {
  name: string;
  slug: string;
  code: string;
}) => {
  const res = await axios.post(`${API_BASE}/admin/scripts`, data);
  return res.data;
};

export const disableScript = async (slug: string) => {
  const res = await axios.post(`${API_BASE}/admin/scripts`, {
    name: "", // required to satisfy server schema
    slug,
    code: "",
    enabled: false,
  });
  return res.data;
};
