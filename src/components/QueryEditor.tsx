import React, { ChangeEvent } from 'react';
import { InlineField, InlineFieldRow, Input } from '@grafana/ui';
import { QueryEditorProps } from '@grafana/data';
import { DataSource } from '../DataSource';
import { PingDataSourceOptions, PingQuery } from '../types';

type Props = QueryEditorProps<DataSource, PingQuery, PingDataSourceOptions>;

export function QueryEditor({ query, onChange }: Props) {
  const onHostChange = (event: ChangeEvent<HTMLInputElement>) => {
    onChange({ ...query, host: event.target.value });
  };

  const onIntervalChange = (event: ChangeEvent<HTMLInputElement>) => {
    onChange({ ...query, interval: event.target.value });
  };

  return (
    <InlineFieldRow>
      <InlineField label="Host" labelWidth={10} tooltip="Hostname or IP address to ping, e.g. 8.8.8.8">
        <Input value={query.host || ''} onChange={onHostChange} placeholder="8.8.8.8" width={24} />
      </InlineField>
      <InlineField
        label="Interval"
        labelWidth={10}
        tooltip="How often to ping this host, e.g. 10s, 30s, 1m"
      >
        <Input value={query.interval || ''} onChange={onIntervalChange} placeholder="10s" width={12} />
      </InlineField>
    </InlineFieldRow>
  );
}
