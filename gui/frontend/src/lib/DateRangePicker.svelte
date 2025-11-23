<script lang="ts">
  import { createEventDispatcher, onMount, onDestroy } from 'svelte';
  import flatpickr from 'flatpickr';
  import 'flatpickr/dist/flatpickr.css';
  import type { Instance } from 'flatpickr/dist/types/instance';
  import type { Options } from 'flatpickr/dist/types/options';

  export let since: Date | null = null;
  export let until: Date | null = null;

  const dispatch = createEventDispatcher();

  let inputElement: HTMLInputElement;
  let flatpickrInstance: Instance | undefined;

  const options: Options = {
    mode: 'range',
    enableTime: true,
    time_24hr: true,
    dateFormat: 'Y-m-d H:i',
    defaultDate: [since, until].filter((d): d is Date => d !== null),
    onChange: (selectedDates: Date[]) => {
      if (selectedDates.length === 2) {
        dispatch('change', { since: selectedDates[0], until: selectedDates[1] });
      } else if (selectedDates.length === 0) {
        dispatch('change', { since: null, until: null });
      }
    },
  };

  onMount(() => {
    if (inputElement) {
      flatpickrInstance = flatpickr(inputElement, options);
    }
  });

  onDestroy(() => {
    flatpickrInstance?.destroy();
  });

  export function clear() {
    flatpickrInstance?.clear();
  }
</script>

<input
  bind:this={inputElement}
  placeholder="Select a date and time range"
  class="form-control"
/>
