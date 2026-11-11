import React, { ChangeEvent } from 'react';
import { InlineField, InlineFieldRow, Input } from '@grafana/ui';
import { DataSourcePluginOptionsEditorProps } from '@grafana/data';
import { PingDataSourceOptions } from '../types';

type Props = DataSourcePluginOptionsEditorProps<PingDataSourceOptions>;

export function ConfigEditor({ options, onOptionsChange }: Props) {
  const jsonData = options.jsonData || {};

  const onIntervalChange = (event: ChangeEvent<HTMLInputElement>) => {
    onOptionsChange({
      ...options,
      jsonData: { ...jsonData, defaultInterval: event.target.value },
    });
  };

  const onTimeoutChange = (event: ChangeEvent<HTMLInputElement>) => {
    onOptionsChange({
      ...options,
      jsonData: { ...jsonData, timeout: event.target.value },
    });
  };

  return (
    <InlineFieldRow>
      <InlineField
        label="Default interval"
        labelWidth={16}
        tooltip="Default ping interval when a query does not specify one, e.g. 10s"
      >
        <Input
          value={jsonData.defaultInterval || ''}
          onChange={onIntervalChange}
          placeholder="10s"
          width={12}
        />
      </InlineField>
      <InlineField
        label="Ping timeout"
        labelWidth={16}
        tooltip="Per-ping timeout, e.g. 3s"
      >
        <Input value={jsonData.timeout || ''} onChange={onTimeoutChange} placeholder="3s" width={12} />
      </InlineField>
    </InlineFieldRow>
  );
}
