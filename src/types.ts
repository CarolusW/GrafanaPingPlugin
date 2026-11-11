import { DataQuery, DataSourceJsonData } from '@grafana/data';

export interface PingQuery extends DataQuery {
  host?: string;
  interval?: string;
}

export interface PingDataSourceOptions extends DataSourceJsonData {
  defaultInterval?: string;
  timeout?: string;
}
