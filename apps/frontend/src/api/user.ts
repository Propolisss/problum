import api from './client';
import type { UserProfile } from '../types';

export async function fetchProfile(): Promise<UserProfile> {
  const resp = await api.get<UserProfile>('/profile');
  return resp.data;
}
