// package routers

// func registerGetValueRouteViaPost(rg *gin.RouterGroup, repository repo.MetricsCRUDer, logger *zap.SugaredLogger) {
// 	rg.POST("/", func(c *gin.Context) {
// 		var metric models.UpdateMetricsModel
// 		var buf bytes.Buffer

// 		// читаем тело запроса
// 		_, err := buf.ReadFrom(c.Request.Body)
// 		if err != nil {
// 			c.JSON(http.StatusBadRequest, gin.H{"error": "failed to read body"})
// 		}

// 		// десериализуем JSON
// 		if err = json.Unmarshal(buf.Bytes(), &metric); err != nil {
// 			c.JSON(http.StatusBadRequest, gin.H{"error": "failed to parse body"})
// 			return
// 		}

// 		if metric.MType == string(models.GaugeType) {
// 			value, err := repository.GetGaugeMetricValueByName(metric.ID, models.GaugeType)
// 			if err == er.ErrorNotFound {
// 				c.JSON(http.StatusNotFound, gin.H{"metric name": metric.ID, "error": "metric not found"})
// 				return
// 			}
// 			c.Header("Content-Type", "application/json")
// 			c.JSON(200, models.UpdateMetricsModel{ID: metric.ID, MType: metric.MType, Value: &value})
// 			return
// 		}
// 		if metric.MType == string(models.CounterType) {
// 			value, err := repository.GetCountMetricValueByName(metric.ID)
// 			if err == er.ErrorNotFound {
// 				c.JSON(http.StatusNotFound, gin.H{"metric name": metric.ID, "error": "metric not found"})
// 				return
// 			}
// 			value64 := int64(value)
// 			c.Header("Content-Type", "application/json")
// 			c.JSON(200, models.UpdateMetricsModel{ID: metric.ID, MType: metric.MType, Delta: &value64})
// 			return
// 		}

// 		c.JSON(http.StatusBadRequest, gin.H{"error": "unknown metric type"})
// 	})
// }

// func registerUpdateRouteViaPost(rg *gin.RouterGroup, repository repo.MetricsCRUDer, syncStorage bool, storagePath string, logger *zap.SugaredLogger) {

// 	rg.POST("/", func(c *gin.Context) {
// 		var metric models.UpdateMetricsModel

// 		var buf bytes.Buffer

// 		// читаем тело запроса
// 		_, err := buf.ReadFrom(c.Request.Body)
// 		if err != nil {
// 			c.JSON(http.StatusBadRequest, gin.H{"error": "failed to read body"})
// 		}

// 		// десериализуем JSON
// 		if err = json.Unmarshal(buf.Bytes(), &metric); err != nil {
// 			c.JSON(http.StatusBadRequest, gin.H{"error": "failed to parse	body"})
// 			return
// 		}

// 		if metric.MType == string(models.GaugeType) {
// 			err := repository.UpdateMetric(metric.ID, models.GaugeType, *metric.Value, syncStorage, storagePath, logger)
// 			if err == er.ErrorNotFound {
// 				c.JSON(http.StatusBadRequest, gin.H{"metric name": metric.ID, "error": "metric not found"})
// 			}
// 			return
// 		}
// 		if metric.MType == string(models.CounterType) {
// 			err := repository.UpdateMetric(metric.ID, models.CounterType, *metric.Delta, syncStorage, storagePath, logger)
// 			if err == er.ErrorNotFound {
// 				c.JSON(http.StatusBadRequest, gin.H{"metric name": metric.ID, "error": "metric not found"})
// 			}
// 			return
// 		}

// 		c.JSON(200, metric)

// 	})

// }

// func registerMultipleUpdateRouteViaPost(rg *gin.RouterGroup, repository repo.MetricsCRUDer, syncStorage bool, storagePath string, logger *zap.SugaredLogger) {

// 	rg.POST("/", func(c *gin.Context) {
// 		var metrics []models.UpdateMetricsModel

// 		var buf bytes.Buffer

// 		// читаем тело запроса
// 		_, err := buf.ReadFrom(c.Request.Body)
// 		if err != nil {
// 			c.JSON(http.StatusBadRequest, gin.H{"error": "failed to read body"})
// 		}

// 		// десериализуем JSON
// 		if err = json.Unmarshal(buf.Bytes(), &metrics); err != nil {
// 			c.JSON(http.StatusBadRequest, gin.H{"error": "failed to parse body"})
// 			return
// 		}

// 		err = repository.UpdateMultipleMetric(metrics)
// 		// TODO: обработать ошибку. вообще все. позаворачивать и тд

// 		if err != nil {
// 			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 		}

// 		c.JSON(200, metrics)

// 	})

// }