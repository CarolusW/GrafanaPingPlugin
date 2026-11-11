import { DataSourcePlugin } from '@grafana/data';
import { PingDataSourceOptions, PingQuery } from './types';
import { DataSource } from './DataSource';
import { ConfigEditor } from './components/ConfigEditor';
import { QueryEditor } from './components/QueryEditor';

export const plugin = new DataSourcePlugin<DataSource, PingQuery, PingDataSourceOptions>(DataSource)
  .setConfigEditor(ConfigEditor)
  .setQueryEditor(QueryEditor);
