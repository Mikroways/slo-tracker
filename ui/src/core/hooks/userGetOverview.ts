import { useEffect, useState } from "react";
import { ISLO } from "../interfaces/ISLO";
import SLOService from "../services/service.slo";

const _sloService = new SLOService();

const useGetOverview = (yearMonth?: string) => {
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [overview, setOverview] = useState<ISLO[]>([])

  const getOverview = async () => {
    setLoading(true);
      setError(null);

      try {
        const response = await _sloService.overview({yearMonth});
        setOverview(response.data.data);
      } catch (err) {
        setError('Error while getting Overview. Please try again');
      } finally {
        setLoading(false);
      }
  }

  useEffect(() => {
    (async() => {
      await getOverview();
    })()
  }, [yearMonth]);

  return {
    loading,
    error,
    overview,
    refreshSLOs: getOverview,
  }
}

export default useGetOverview;
