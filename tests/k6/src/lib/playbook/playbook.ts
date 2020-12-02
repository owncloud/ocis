import {Gauge, Trend} from "k6/metrics";

export class Play {
    public readonly name: string;
    public readonly metricTrendName: string;
    public readonly metricErrorRateName: string;
    public readonly metricTrend: Trend;
    public readonly metricErrorRate: Gauge;
    protected tags: { [key: string]: string };

    constructor({name}: { name: string; }) {
        this.name = name;
        this.metricTrendName = `${this.name}_trend`;
        this.metricErrorRateName = `${this.name}_error_rate`;
        this.metricTrend = new Trend(this.metricTrendName, true);
        this.metricErrorRate = new Gauge(this.metricErrorRateName);
        this.tags = {play: this.name}
    }
}