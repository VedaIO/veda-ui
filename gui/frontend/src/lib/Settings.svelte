<script lang="ts">
  import { onMount } from 'svelte';
  import { writable } from 'svelte/store';
  import { openUninstallModal } from './modalStore';

  let isAutostartEnabled = false;
  let autostartStatusText = writable('Không rõ');
  let autostartStatusClass = writable('bg-secondary');
  let autostartToggleBtnText = writable('Toggle');
  let autostartToggleBtnDisabled = false;

  async function loadAutostartStatus(): Promise<void> {
    try {
      const res = await fetch('/api/settings/autostart/status');
      if (!res.ok) {
        autostartStatusText.set('Không hỗ trợ trên HĐH này');
        autostartStatusClass.set('bg-warning');
        autostartToggleBtnDisabled = true;
        return;
      }
      const data = await res.json();
      isAutostartEnabled = data.enabled;

      autostartStatusText.set(isAutostartEnabled ? 'Đã bật' : 'Đã tắt');
      autostartStatusClass.set(
        isAutostartEnabled ? 'bg-success' : 'bg-secondary'
      );

      autostartToggleBtnText.set(
        isAutostartEnabled ? 'Tắt tự động khởi động' : 'Bật tự động khởi động'
      );
      autostartToggleBtnDisabled = false;
    } catch {
      autostartStatusText.set('Lỗi');
      autostartStatusClass.set('bg-danger');
      autostartToggleBtnDisabled = true;
    }
  }

  function toggleAutostart(): void {
    autostartToggleBtnDisabled = true;
    const endpoint = isAutostartEnabled
      ? '/api/settings/autostart/disable'
      : '/api/settings/autostart/enable';
    try {
      fetch(endpoint, { method: 'POST' }).then((res) => {
        if (!res.ok) {
          res.text().then((errorText) => {
            alert(`Thao tác thất bại: ${errorText}`);
          });
        } else {
          loadAutostartStatus(); // Refresh status after action
        }
      });
    } catch (e) {
      if (e instanceof Error) {
        alert(`Đã xảy ra lỗi: ${e.message}`);
      }
    } finally {
      autostartToggleBtnDisabled = false;
    }
  }

  onMount(() => {
    loadAutostartStatus();
  });
</script>

<div id="settings-view">
  <h2 class="mb-4">Cài đặt</h2>

  <div class="row">
    <div class="col-md-8 mx-auto">
      <!-- Autostart Settings Card -->
      <div class="card mb-4">
        <div class="card-header">
          <h4>Khởi động cùng Windows</h4>
        </div>
        <div class="card-body">
          <p class="card-text">
            Trạng thái:
            <span class="badge {$autostartStatusClass}"
              >{$autostartStatusText}</span
            >
          </p>
          <p class="card-text">
            Bật tùy chọn này để ProcGuard tự động chạy khi bạn đăng nhập vào
            Windows.
          </p>
          <button
            id="autostart-toggle-btn"
            class="btn btn-primary"
            on:click={toggleAutostart}
            disabled={autostartToggleBtnDisabled}
          >
            {$autostartToggleBtnText}
          </button>
        </div>
      </div>

      <!-- Uninstall Card -->
      <div class="card mb-4">
        <div class="card-header">
          <h4>Gỡ cài đặt ProcGuard</h4>
        </div>
        <div class="card-body">
          <p class="card-text">
            <b>Cảnh báo:</b> Thao tác này sẽ xóa toàn bộ dữ liệu và gỡ cài đặt ProcGuard
            khỏi hệ thống.
          </p>
          <button
            type="button"
            class="btn btn-danger"
            on:click={openUninstallModal}
          >
            Gỡ cài đặt ProcGuard
          </button>
        </div>
      </div>
    </div>
  </div>
</div>
