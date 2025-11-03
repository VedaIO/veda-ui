import { writable } from 'svelte/store';

export const isUninstallModalOpen = writable(false);
export const uninstallPassword = writable('');
export const uninstallError = writable('');

export function openUninstallModal() {
  uninstallPassword.set('');
  uninstallError.set('');
  isUninstallModalOpen.set(true);
}

export async function handleUninstallSubmit() {
  let password = '';
  uninstallPassword.subscribe((value) => (password = value))();

  const response = await fetch('/api/uninstall', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ password }),
  });

  if (response.ok) {
    isUninstallModalOpen.set(false);
    alert('ProcGuard đã được gỡ cài đặt. Trang sẽ đóng.');
    window.location.href = 'about:blank';
  } else {
    uninstallError.set('Gỡ cài đặt thất bại. Vui lòng kiểm tra lại mật khẩu.');
  }
}
