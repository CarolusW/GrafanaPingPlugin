import { DataQueryRequest, DataFrame, DataSourceApi, DataSourceInstanceSettings } from '@grafana/data';
import { PingDataSourceOptions, PingQuery } from './types';

export class DataSource extends DataSourceApi<PingQuery, PingDataSourceOptions> {
  constructor(instanceSettings: DataSourceInstanceSettings<PingDataSourceOptions>) {
    super(instanceSettings);
  }

  async query(request: DataQueryRequest<PingQuery>): Promise<{ data: DataFrame[] }> {
    return { data: [] };
  }

  async testDatasource() {
    return { status: 'success', message: 'Ping datasource configured' };
  }
}
