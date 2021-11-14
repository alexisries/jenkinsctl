package jobs

import (
	"fmt"
	"jenkinsctl/pkg/apiclient"
	"regexp"

	"github.com/beevik/etree"
)

func selectOrCreateElement(root *etree.Element, name string) *etree.Element {
	element := root.SelectElement(name)
	if element == nil {
		element = root.CreateElement(name)
	}
	return element
}

func (jobs *Jobs) Schedule(clt *apiclient.ApiClient, schedule string) error {
	fmt.Println("Scheduling jobs...")

	var xmlHeaderRegex = regexp.MustCompile(`^<\?xml.*\?>`)

	for _, job := range jobs.Jobs {
		rawXml, _ := job.JenkinsJob.GetConfig(clt.Ctx)
		rawXml = xmlHeaderRegex.ReplaceAllString(rawXml, "")

		doc := etree.NewDocument()
		if err := doc.ReadFromString(rawXml); err != nil {
			return err
		}
		flowDefinition := doc.SelectElement("flow-definition")
		properties := selectOrCreateElement(flowDefinition, "properties")
		pipelineTriggersJobProperty := selectOrCreateElement(
			properties, "org.jenkinsci.plugins.workflow.job.properties.PipelineTriggersJobProperty",
		)
		triggers := selectOrCreateElement(pipelineTriggersJobProperty, "triggers")
		timerTrigger := selectOrCreateElement(triggers, "hudson.triggers.TimerTrigger")
		spec := selectOrCreateElement(timerTrigger, "spec")

		before := spec.Text()
		if before == schedule {
			fmt.Printf("%s schedule is already defined on %s job\n", schedule, job.Name)
			continue
		}
		spec.SetText(schedule)

		doc.Indent(2)
		newXml, err := doc.WriteToString()
		if err != nil {
			return err
		}

		err = job.JenkinsJob.UpdateConfig(clt.Ctx, newXml)
		if err != nil {
			return err
		}
		fmt.Printf("job %s is now scheduled (%s)\n", job.Name, schedule)
	}
	return nil
}
